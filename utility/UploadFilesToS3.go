package utility

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
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

func parseUploadFilesToS3Request(c *gin.Context) (*UploadFilesToS3Request, error) {
	// PUT body params
	var ufr UploadFilesToS3Request

	if err := c.ShouldBindJSON(&ufr); err != nil {
		return nil, err
	}

	return &ufr, nil
}

func UploadFilesToS3(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	uploadFileRequest, err := parseUploadFilesToS3Request(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// validation //TODO
	for source, urls := range uploadFileRequest.Source {

		if source != "gdrive" { // valid urls are gdrive
			utils.GeneralBadRequestError(c, "Invalid source")
			return
		}

		if len(urls) > 10 {
			utils.GeneralBadRequestError(c, "Maximum 10 files can be uploaded at a time")
			return
		}
	}

	// check if community is not on free plan
	if validateIfCommunityIsOnFreeTier(headers["x-api-key"].(string)) {
		utils.GeneralBadRequestError(c, "please upgrade your tier plan to upload data to S3 directly")
		return
	}

	// initalise file_path
	file_path := ""

	// switch case for file path
	if uploadFileRequest.Category == "feed" {
		switch uploadFileRequest.Entity {
		case "post":
			file_path = fmt.Sprintf("files/post/%s", userId)
		case "widget":
			file_path = fmt.Sprintf("files/widget/%s", userId)
		default:
			utils.GeneralBadRequestError(c, "Invalid entity")
		}
	} else {
		utils.GeneralBadRequestError(c, "Invalid category")
		return
	}

	// file upload response
	fileUploadedUrls := []FileUploadedResponse{}

	for source, urls := range uploadFileRequest.Source {

		if source == "gdrive" { //constatns
			filesUploaded := uploadFilesFromDrive(urls, file_path)
			fileUploadedUrls = append(fileUploadedUrls, filesUploaded...)
		}
	}

	response := gin.H{
		"files": fileUploadedUrls,
	}

	utils.GenerateResponse(c, response, false)
}

// upload files from drive to s3 in parallel
func uploadFilesFromDrive(urls []string, file_path string) []FileUploadedResponse {

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
				s3Url, err := uploadDriveFileToS3(fileID, file_path)
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

	// regex pattern
	// pattern := `\/(?:file\/d\/|open\?id=|uc\?id=)([a-zA-Z0-9_-]+)`
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

func uploadDriveFileToS3(fileId string, filePath string) (string, error) {

	// Call lambda function to upload file to s3
	lambdaUrl := "https://kqmg7qu7uzgcnqkxmrjafi5n3q0rbrgb.lambda-url.ap-south-1.on.aws" //TODO: constants
	headers := gin.H{
		"x-platform-type": "kettle-service",
		"content-type":    "application/json",
	}
	body := gin.H{
		"file_id":   fileId,
		"file_path": filePath,
	}

	respBytes, _, err := utils.CallExternalAPI(lambdaUrl, utils.PUTRequest, headers, nil, body)
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
