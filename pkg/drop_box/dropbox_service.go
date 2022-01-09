package drop_box

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type DropboxService struct {
	Token  string
	Client files.Client
	Id     bool
}

func NewDropboxService(token string, options...bool) (*DropboxService, error) {
	config := dropbox.Config{
		Token: token,
	}
	client := files.New(config)

	id := false
	if len(options) > 0 {
		id = options[0]
	}
	return &DropboxService{Token: token, Client: client, Id: id}, nil
}

func (d DropboxService) Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error) {
	file := bytes.NewReader(data)

	// create new client to access drop_box cloud with token generated in drop_box console
	client := d.Client
	if client == nil {
		config := dropbox.Config{Token: d.Token}
		client = files.New(config)
	}

	// create new upload info
	var filepath string
	if len(directory) > 0 {
		filepath = fmt.Sprintf("/%s/%s", directory, filename)
	} else {
		filepath = "/" + filename
	}

	arg := files.NewCommitInfo(filepath)

	// upload file
	res, err := client.Upload(arg, file)
	if err != nil {
		return "", nil
	}
	// msg := fmt.Sprintf("uploaded file '%s' to dropbox successfully!!!", filename)
	if d.Id {
		return res.Id, nil
	} else {
		return res.PathLower + "/" + res.Name, nil
	}
}

func (d DropboxService) Delete(ctx context.Context, directory string, filename string) (bool, error) {
	client := d.Client
	if client == nil {
		config := dropbox.Config{Token: d.Token}
		client = files.New(config)
	}
	var filepath string
	if len(directory) > 0 {
		filepath = fmt.Sprintf("/%s/%s", directory, filename)
	} else {
		filepath = "/" + filename
	}
	arg := files.NewDeleteArg(filepath)
	_, err := client.DeleteV2(arg)
	if err != nil {
		return false, err
	}
	return true, nil
}
