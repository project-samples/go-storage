package drop_box

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type DropboxService struct {
	Token	string
	Client	files.Client
}

func NewDropboxService(token string) (*DropboxService, error) {
	config := dropbox.Config{
		Token: token,
	}
	client := files.New(config)

	return &DropboxService{Token: token, Client: client}, nil
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

	// upload file
	_, err := client.Upload(arg, file)
	if err != nil {
		panic(err)
	}

	msg := fmt.Sprintf("uploaded file '%s' to dropbox successfully!!!", filename)
	return msg, err
}

func (d DropboxService)  Delete(ctx context.Context, directory string, fileName string) (bool, error) {
	client := d.Client
	if client == nil {
		config := dropbox.Config{
			Token: d.Token,
		}
		client = files.New(config)
	}

	filepath := fmt.Sprintf("/%s/%s", directory, fileName)
	arg := files.NewDeleteArg(filepath)
	_, err := client.DeleteV2(arg)
	if err != nil {
		return false, err
	}

	return true, nil
}