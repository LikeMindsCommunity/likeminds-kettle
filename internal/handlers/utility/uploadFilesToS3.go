package utility

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/environment"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// PUT Request structure for uploading files to S3
type UploadFilesToS3Request struct {
	Category string              `json:"category" binding:"required"`
	Entity   string              `json:"entity" binding:"required"`
	Source   map[string][]string `json:"source" binding:"required"`
}

// Response structure for file upload serverless function
type FileUploadedResponse struct {
	SourceUrl string `json:"source_url"`
	S3Url     string `json:"s3_url"`
	Error     string `json:"error"`
}

func parseAndValidateUploadFilesToS3Request(c *gin.Context, headers map[string]interface{}) (string, map[string][]string, error) {
	var ufr UploadFilesToS3Request
	if err := c.ShouldBindJSON(&ufr); err != nil {
		return "", nil, err
	}

	for source, urls := range ufr.Source {

		if source != SourceGDrive {
			return "", nil, fmt.Errorf(utils.ErrorInvalidSource)
		}

		if len(urls) > MaxFilesPerUpload {
			return "", nil, fmt.Errorf(utils.ErrorMaxFilesPerUpload)
		}
	}

	// validate community tier
	if utils.IsCommunityOnFreeTier(utils.GetRedisClientFromContext(c), headers) {
		return "", nil, fmt.Errorf(utils.ErrorUpgradeTierPlanForS3)
	}

	var filePath string
	switch ufr.Category {
	case CategoryFeed:
		switch ufr.Entity {
		case EntityPost:
			filePath = fmt.Sprintf(FeedPostFilePath, headers[utils.HeadersMemberId].(string))
		case EntityWidget:
			filePath = fmt.Sprintf(FeedWidgetFilePath, headers[utils.HeadersMemberId].(string))
		default:
			return "", nil, fmt.Errorf(utils.ErrorInvalidEntity)
		}
	default:
		return "", nil, fmt.Errorf(utils.ErrorInvalidCategory)
	}

	return filePath, ufr.Source, nil
}

func UploadFilesToS3(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	filePath, sourceUrls, err := parseAndValidateUploadFilesToS3Request(c, headers)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// file upload response
	fileUploadedUrls := []FileUploadedResponse{}

	for source, urls := range sourceUrls {
		if source == SourceGDrive {
			filesUploaded := uploadFilesFromDrive(urls, filePath)
			fileUploadedUrls = append(fileUploadedUrls, filesUploaded...)
		}
	}

	response := gin.H{
		"files": fileUploadedUrls,
	}

	utils.GenerateResponse(c, response, false)
}

// upload files from drive to s3 in parallel
func uploadFilesFromDrive(urls []string, filePath string) []FileUploadedResponse {

	var wg sync.WaitGroup
	fileUploadedUrls := make([]FileUploadedResponse, len(urls))

	for i, url := range urls {
		wg.Add(1)
		utils.SafeGo(func(i int, url string) func() {
			return func() {
				defer wg.Done()
				fileResponse := FileUploadedResponse{
					SourceUrl: url,
				}

				// extract file id from url
				fileID, err := extractFileIdFromDriveUrl(url)
				if err != nil {
					fileResponse.Error = err.Error()
				} else {
					// upload file to s3
					s3Url, err := uploadDriveFileToS3(fileID, filePath)
					if err != nil {
						fileResponse.Error = err.Error()
					} else {
						fileResponse.S3Url = s3Url
					}
				}
				fileUploadedUrls[i] = fileResponse
			}
		}(i, url))
	}

	wg.Wait()
	return fileUploadedUrls
}

// extracts file_id from gdrive share or downloadin link using regex
func extractFileIdFromDriveUrl(url string) (string, error) {

	// Regex pattern
	regex, err := regexp.Compile(FilterFileIDFromDriveUrlRegex)
	if err != nil {
		return "", err
	}

	// Find the match in the link
	match := regex.FindStringSubmatch(url)
	if match == nil || len(match) < 2 {
		return "", fmt.Errorf("no file ID found in the link")
	}

	// Extract the file ID
	fileID := match[1]
	return fileID, nil
}

// upload google drive file to s3 using serverless function
func uploadDriveFileToS3(fileId string, filePath string) (string, error) {

	// Call serverless function to upload file to s3
	functionUrl := environment.GoDotEnvVariable(environment.EnvUploadFromDriveToS3LambdaUrl)

	headers := gin.H{
		utils.HeadersPlatformType: string(utils.PlatformKettleService),
		utils.HeaderContentType:   utils.ContentTypeApplicationJson,
	}

	// check if beta environment
	isBeta := true
	serverEnv := environment.GoDotEnvVariable(environment.EnvServerEnviornment)
	if serverEnv == "prod" {
		isBeta = false
	}

	body := gin.H{
		"is_prod":   !isBeta,
		"file_id":   fileId,
		"file_path": filePath,
	}

	respBytes, _, err := utils.CallExternalAPI(functionUrl, utils.PUTRequest, headers, nil, body)
	if err != nil {
		return "", fmt.Errorf("some error occured while uploading file to s3")
	}

	fileUploadResponse := FileUploadedResponse{}
	if err := json.Unmarshal(respBytes, &fileUploadResponse); err != nil {
		return "", err
	}

	// if error in response
	if fileUploadResponse.Error != "" {
		return "", fmt.Errorf(fileUploadResponse.Error)
	}

	// return s3 url
	return fileUploadResponse.S3Url, nil
}
