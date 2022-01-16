package drop_box

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"io/ioutil"
	"net/http"
)

type DropboxService struct {
	Token     string
	Client    files.Client
	ClientAPI *http.Client
	Id        bool
}

type BodyRequestDropbox struct {
	Path     string         `json:"path"`
	Settings SettingsStruct `json:"settings"`
}

type SettingsStruct struct {
	Audience            string `json:"audience"`
	Access              string `json:"access"`
	RequestedVisibility string `json:"requested_visibility"`
	AllowDownload       bool   `json:"allow_download"`
}

type FileShareResponse struct {
	Tag             string `json:".tag"`
	Url             string `json:"url"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	PathLower       string `json:"path_lower"`
	LinkPermissions string `json:",omitempty"`
	// linkPermissions has a very complicated and large structure data.
	// Because its data is not important in this project, so it will be ignored.
	PreviewType    string  `json:"preview_type"`
	ClientModified string  `json:"client_modified"`
	ServerModified string  `json:"server_modified"`
	Rev            string  `json:"rev"`
	Size           float64 `json:"size"`
}

func NewDropboxService(token string) (*DropboxService, error) {
	config := dropbox.Config{
		Token: token,
	}
	client := files.New(config)
	clientAPI := &http.Client{}
	return &DropboxService{Token: token, Client: client, ClientAPI: clientAPI}, nil
}

func (d DropboxService) Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error) {
	file := bytes.NewReader(data)

	// create new client to access drop_box cloud with token generated in drop_box console
	client := d.Client
	if client == nil {
		config := dropbox.Config{
			Token: d.Token,
		}
		client = files.New(config)
	}

	// create new upload info
	filepath := fmt.Sprintf("/%s/%s", directory, filename)
	arg := files.NewCommitInfo(filepath)

	//upload file
	_, err := client.Upload(arg, file)
	if err != nil {
		return "", err
	}

	//consuming API
	url := "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"
	var bearer = "Bearer " + d.Token
	bodyData := BodyRequestDropbox{
		Path: filepath,
		Settings: SettingsStruct{
			Audience:            "public",
			Access:              "viewer",
			RequestedVisibility: "public",
			AllowDownload:       true,
		},
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(bodyData)
	if err != nil {
		return "", err
	}

	//bodyData := map[string]interface{}{
	//	"path": filepath,
	//	"settings":  map[string]interface{}{
	//		"audience": "public",
	//		"access": "viewer",
	//		"requested_visibility": "public",
	//		"allow_download": true,
	//	},
	//}
	//byteArray, _ := json.Marshal(bodyData)
	//reader := bytes.NewReader(byteArray)

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")
	// Send req using http Client
	resp, err := d.ClientAPI.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var bodyResp FileShareResponse
	err = json.Unmarshal(body, &bodyResp)
	if err != nil {
		return "", err
	}
	if d.Id {
		return bodyResp.Id, nil
	} else {
		return bodyResp.Url, nil
	}
	// msg := fmt.Sprintf("uploaded file '%s' to dropbox successfully!!! follow this link to view file in dropbox: %s", filename, bodyResp.Url)
	// return msg, err
}

func (d DropboxService) Delete(ctx context.Context, fileName string) (bool, error) {
	client := d.Client
	if client == nil {
		config := dropbox.Config{
			Token: d.Token,
		}
		client = files.New(config)
	}
	// filepath := fmt.Sprintf("/%s/%s", directory, fileName)
	arg := files.NewDeleteArg(fileName)
	_, err := client.DeleteV2(arg)
	if err != nil {
		return false, err
	}
	return true, nil
}
