package one_drive

import (
	"context"
	"fmt"
	"github.com/goh-chunlin/go-onedrive/onedrive"
	"golang.org/x/oauth2"
	"log"
	"os"
)

type OneDriveService struct {
	AccessToken	string
	Client *onedrive.Client
}

func NewOneDriveService(ctx context.Context, token string) (*OneDriveService, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := onedrive.NewClient(tc)

	return &OneDriveService{AccessToken: token, Client: client}, nil
}

func (o OneDriveService) Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error) {
	// code here
	client := o.Client
	if client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: o.AccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = onedrive.NewClient(tc)
	}

	// create a folder to store file to upload
	uploadFolder := "one_drive_upload_file"
	isExists, err := exists(uploadFolder)
	if isExists == false {
		err = os.Mkdir(uploadFolder, 0755)
		if err!= nil {
			panic(err)
		}
	}

	// create file to upload
	filePath := fmt.Sprintf("%s/%s", uploadFolder, filename)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	_, err2 := file.Write(data)
	if err2 != nil {
		panic(err2)
	}

	item, err := client.DriveItems.UploadNewFile(ctx, "", "root", filePath)
	if err != nil {
		return "", err
	}

	msg := fmt.Sprintf("uploaded file '%s' to one-drive successfully!!!, follow this link to view file in one-drive: %s", filename, item.WebURL)
	return msg, nil
}

func (o OneDriveService) Delete(ctx context.Context, directory string, fileName string) (bool, error) {
	client := o.Client
	if client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: o.AccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = onedrive.NewClient(tc)
	}

	list,_ := client.DriveItems.List(ctx, "root")
	var fileId string
	for i := 0; i < len(list.DriveItems); i++ {
		if list.DriveItems[i].Name == fileName {
			fileId = list.DriveItems[i].Id
		}
	}

	driveResp, _ := client.Drives.List(ctx)
	var driveId string
	for i := 0; i < len(driveResp.Drives); i++ {
		if driveResp.Drives[i].DriveType == "personal" {
			driveId = driveResp.Drives[i].Id
		}
	}
	log.Print(driveId)

	err := client.DriveItems.Delete(ctx, "", fileId)
	if err != nil {
		return false, err
	}

	return true, err
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}