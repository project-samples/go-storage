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
	Token  string
	Client *onedrive.Client
	Id     bool
}

func NewOneDriveService(ctx context.Context, token string, options...bool) (*OneDriveService, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := onedrive.NewClient(tc)
	id := true
	if len(options) > 0 {
		id = options[0]
	}
	return &OneDriveService{Token: token, Client: client, Id: id}, nil
}

func (o OneDriveService) Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error) {
	// code here
	client := o.Client
	if client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: o.Token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = onedrive.NewClient(tc)
	}

	// create a folder to store file to upload
	uploadFolder := "one_drive_upload_file"
	isExists, err := Exists(uploadFolder)
	if isExists == false {
		err = os.Mkdir(uploadFolder, 0755)
		if err != nil {
			return "", err
		}
	}

	// create file to upload
	filePath := fmt.Sprintf("%s/%s", uploadFolder, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	_, err2 := file.Write(data)
	if err2 != nil {
		return "", err2
	}

	item, err := client.DriveItems.UploadNewFile(ctx, "", "root", filePath)
	if err != nil {
		return "", err
	}
	if o.Id {
		return item.Id, nil
	} else {
		return item.WebURL, nil
	}
	// msg := fmt.Sprintf("uploaded file '%s' to one-drive successfully!!!, follow this link to view file in one-drive: %s", filename, item.WebURL)
	// return msg, nil
}

func (o OneDriveService) Delete(ctx context.Context, fileName string) (bool, error) {
	client := o.Client
	if client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: o.Token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = onedrive.NewClient(tc)
	}
	if o.Id {
		err := client.DriveItems.Delete(ctx, "", fileName)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		list, _ := client.DriveItems.List(ctx, "root")
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
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
