package utils

import (
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"posthis/entity"

	uuid "github.com/satori/go.uuid"
)

type Media = entity.Media

func UploadMultipleFiles(files []*multipart.FileHeader, media *[]*Media) error {

	for i, _ := range files {
		file, err := files[i].Open()
		if err != nil {
			return err
		}
		defer file.Close()

		name := uuid.NewV4().String() + files[i].Filename

		bytes := make([]byte, files[i].Size)

		_, err = file.Read(bytes)
		if err != nil {
			return err
		}

		mime := http.DetectContentType(bytes)

		err = ioutil.WriteFile(filepath.Join("static", name), bytes, fs.ModePerm)
		if err != nil {
			return err
		}

		*media = append(*media, &Media{Name: name, Mime: mime})
	}

	return nil
}
