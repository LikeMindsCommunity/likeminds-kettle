package utility

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/environment"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// PUT Request structure for uploading files to S3
type UploadFilesToS3Request struct {
	Category string              `json:"category"`
	Entity   string              `json:"entity"`
	Source   map[string][]string `json:"source"`
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

		if source != "gdrive" { // valid urls are gdrive
			return "", nil, fmt.Errorf("invalid source")
		}

		if len(urls) > 10 {
			return "", nil, fmt.Errorf("maximum 10 files can be uploaded at a time")
		}
	}

	// check if community is not on free plan
	if validateIfCommunityIsOnFreeTier(headers[utils.HeadersApiKey].(string)) {
		return "", nil, fmt.Errorf("please upgrade your tier plan to upload data to S3 directly")
	}

	var filePath string

	if ufr.Category == CategoryFeed {
		switch ufr.Entity {
		case EntityPost:
			filePath = fmt.Sprintf("files/post/%s", headers[utils.HeadersMemberId].(string))
		case EntityWidget:
			filePath = fmt.Sprintf("files/widget/%s", headers[utils.HeadersMemberId].(string))
		default:
			return "", nil, fmt.Errorf("invalid entity")
		}
	} else {
		return "", nil, fmt.Errorf("invalid category")
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
		go func(i int, url string) {
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
		}(i, url)
	}

	wg.Wait()
	return fileUploadedUrls
}

// extracts file_id from gdrive share or downloadin link using regex
func extractFileIdFromDriveUrl(url string) (string, error) {

	// Regex pattern
	pattern := `\/(?:file\/d\/|open\?id=|uc\?export=download&id=)([a-zA-Z0-9_-]+)`
	regex, err := regexp.Compile(pattern)
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
	functionUrl := environment.GoDotEnvVariable("UPLOAD_DRIVE_TO_S3_FUNCTION_URL")
	headers := gin.H{
		"x-platform-type": "kettle-service",
		"content-type":    "application/json",
	}
	body := gin.H{
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

func validateIfCommunityIsOnFreeTier(_ string) bool {
	return false // TODO
}
