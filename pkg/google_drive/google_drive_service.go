package google_drive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io"
	"log"
	"net/http"
	"os"
)

type GoogleDriveService struct {
	Service *drive.Service
	Id      bool
}

func NewGoogleDriveService(credentials []byte, options ...bool) (*GoogleDriveService, error) {
	// create a config struct based on Google Drive credentials string in config
	config, err := google.ConfigFromJSON(credentials, drive.DriveFileScope)
	if err != nil {
		return nil, err
	}

	client := GetClient(config)

	service, err := drive.New(client)
	if err != nil {
		fmt.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}
	id := true
	if len(options) > 0 {
		id = options[0]
	}
	return &GoogleDriveService{Service: service, Id: id}, err
}

func (s GoogleDriveService) Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error) {
	file := bytes.NewReader(data)
	var fileId string
	var folderId string

	//check duplicate folder
	queryFolder := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.folder' and trashed = false", directory)
	listFolder, _ := s.Service.Files.List().Q(queryFolder).Do()
	if listFolder != nil && len(listFolder.Files) > 0 {
		folderId = listFolder.Files[0].Id
	} else {
		folder, err1 := CreateDirectory(s.Service, directory, "root")
		if err1 != nil {
			return "could not create directory!!!", err1
		}
		folderId = folder.Id
	}

	// check duplicate file
	queryFile := fmt.Sprintf("name = '%s' and mimeType != 'application/vnd.google-apps.folder' and trashed = false", filename)
	listFile, _ := s.Service.Files.List().Q(queryFile).Do()
	if listFile != nil && len(listFile.Files) > 0 {
		fileId = listFile.Files[0].Id
		err := s.Service.Files.Delete(fileId).Do()
		if err != nil {
			return "could not delete duplicate file", err
		}
	}

	// create the file and upload its content
	fileResp, err := CreateFile(s.Service, filename, contentType, file, folderId)
	if err != nil {
		msg := fmt.Sprintf("Could not create file: %v\n", err)
		return msg, err
	}
	fileId = fileResp.Id

	// 		create share link for file
	// create permission type
	permissionToBeCreated := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	// create permission for file
	_, err = s.Service.Permissions.Create(fileId, permissionToBeCreated).Do() // create permission for file
	if err != nil {
		return fileId, err
	}
	if s.Id {
		return fileResp.Id, nil
	} else {
		//get webViewLink
		fileRespShare, err := s.Service.Files.Get(fileId).Fields("webViewLink").Do() //get webViewLink
		if err != nil {
			return "", err
		}
		return fileRespShare.WebViewLink, nil
	}
}

func (s GoogleDriveService) Delete(ctx context.Context, fileId string) (bool, error) {
	if s.Id {
		err := s.Service.Files.Delete(fileId).Do()
		if err != nil {
			return false, err
		} else {
			return true, err
		}
	} else {
		// get the fileId of the file that need to be deleted
		q := fmt.Sprintf("name = '%s' and mimeType != 'application/vnd.google-apps.folder' and trashed = false", fileId)
		list, err := s.Service.Files.List().Q(q).Do()
		if list == nil || err != nil {
			return false, err
		}
		fileId := list.Files[0].Id

		err = s.Service.Files.Delete(fileId).Do()
		if err != nil {
			return false, err
		}
		return true, err
	}
}

func CreateDirectory(service *drive.Service, name string, parentId string) (*drive.File, error) {
	// common parentId is "root"
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

func CreateFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

// getClient Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := TokenFromFile(tokFile)

	// token file not exists, get token from web based on config, then save token info in the new token file
	if err != nil {
		tok = GetTokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// getTokenFromWeb Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenFromFile get token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken Saves a token to a new created file
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
