package dropboxcli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/sharing"
	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.uber.org/zap"
)

const (
	configExpirationDuration = 30 * time.Minute
)

type storageDropbox struct {
	appKey               string
	appSecret            string
	refreshToken         string
	httpClient           *http.Client
	dropBoxConfig        *dropbox.Config
	lastConfigUpdateTime *time.Time
	log                  *zap.Logger
}

func NewStorageDropbox(
	appKey string,
	appSecret string,
	refreshToken string,
	httpClient *http.Client,
	log *zap.Logger) repository.StorageRepository {
	return &storageDropbox{
		appKey:        appKey,
		appSecret:     appSecret,
		refreshToken:  refreshToken,
		httpClient:    httpClient,
		dropBoxConfig: nil,
		log:           log,
	}
}

func (port *storageDropbox) GetStorageType() string {
	return "DROPBOX"
}

func (port *storageDropbox) GetFileUploadUrl(metaData domain.FileMetaData) (string, error) {

	filePath := "/" + metaData.Path

	port.log.Debug("GetFileUploadUrl inputs: ",
		zap.String("filePath", filePath))

	filesClient, err := port.getFilesClient()
	if err != nil {
		return "", err
	}

	uploadLinkResponse, err := filesClient.GetTemporaryUploadLink(files.NewGetTemporaryUploadLinkArg(files.NewCommitInfo(filePath)))
	if err != nil {
		return "", err
	}

	port.log.Debug("temp upload link response: ", zap.Any("data", uploadLinkResponse))

	return uploadLinkResponse.Link, nil
}

func (port *storageDropbox) GetFileDownloadUrl(metaData domain.FileMetaData) (string, error) {

	filePath := "/" + metaData.Path

	port.log.Debug("GetFileDownloadUrl inputs: ",
		zap.String("filePath", filePath))

	sharingClient, err := port.getSharingClient()
	if err != nil {
		return "", nil
	}

	linkMetaData, err := sharingClient.CreateSharedLink(sharing.NewCreateSharedLinkArg(filePath))
	if err != nil {
		return "", err
	}

	port.log.Debug("share download link response: ", zap.Any("data", linkMetaData))

	return port.cleanDownloadUrl(linkMetaData.Url), nil
}

func (port *storageDropbox) DeleteFile(metaData domain.FileMetaData) error {

	filePath := "/" + metaData.Path

	port.log.Debug("DeleteFile inputs: ",
		zap.String("filePath", filePath))

	filesClient, err := port.getFilesClient()
	if err != nil {
		return err
	}

	deleteResponse, err := filesClient.DeleteV2(files.NewDeleteArg(filePath))
	if err != nil {
		return err
	}

	port.log.Debug("delete response: ", zap.Any("data", deleteResponse))

	return nil
}

// private

func (port *storageDropbox) getFilesClient() (files.Client, error) {
	config, err := port.getDropBoxConfig()
	if err != nil {
		return nil, err
	}
	return files.New(*config), nil
}

func (port *storageDropbox) getSharingClient() (sharing.Client, error) {
	config, err := port.getDropBoxConfig()
	if err != nil {
		return nil, err
	}
	return sharing.New(*config), nil
}

func (port *storageDropbox) getDropBoxConfig() (*dropbox.Config, error) {

	configExpired := port.lastConfigUpdateTime != nil && time.Now().After(port.lastConfigUpdateTime.Add(configExpirationDuration))
	if port.dropBoxConfig != nil && !configExpired {
		return port.dropBoxConfig, nil
	}

	accessToken, err := port.getAccessToken(
		port.appKey,
		port.appSecret,
		port.refreshToken,
		port.httpClient,
	)
	if err != nil {
		return nil, err
	}

	port.log.Debug("accesToken", zap.String("value", accessToken))

	config := dropbox.Config{
		Token:    accessToken,
		LogLevel: dropbox.LogInfo,
	}

	port.dropBoxConfig = &config
	nowTime := time.Now()
	port.lastConfigUpdateTime = &nowTime
	return &config, nil
}

func (port *storageDropbox) getAccessToken(
	appKey string,
	appSecret string,
	refreshToken string,
	httpClient *http.Client) (accessToken string, outputErr error) {

	dropboxTokenUrl := "https://api.dropbox.com/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	request, err := http.NewRequest(
		http.MethodPost,
		dropboxTokenUrl,
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	request.SetBasicAuth(appKey, appSecret)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := port.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer func() {
		outputErr = response.Body.Close()
	}()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	var jsonResult map[string]any
	err = json.NewDecoder(response.Body).Decode(&jsonResult)
	if err != nil {
		return "", err
	}

	return jsonResult["access_token"].(string), nil
}

func (port *storageDropbox) cleanDownloadUrl(downloadUrl string) string {
	return strings.Replace(downloadUrl, "?dl=0", "?raw=1", 1)
}
