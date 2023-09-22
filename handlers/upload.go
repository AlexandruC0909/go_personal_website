package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UploadHandler interface {
	UploadImages(w http.ResponseWriter, r *http.Request) error
}

// handler to handle the image upload
func (s *ApiRouter) UploadImages(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)
	if err != nil {
		return err
	}

	user, err := s.store.GetUser(id)
	if err != nil {
		return err
	}
	if err := r.ParseMultipartForm(32 * 1024 * 1024); err != nil {
		return err
	}

	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			return err
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" {
			return fmt.Errorf("The provided file format is not allowed. Please upload a JPEG, JPG, or PNG image")
		}

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			return err
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		newFileName := fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		f, err := os.Create(newFileName)
		if err != nil {
			return err
		}
		user.ImageURL = newFileName

		s.store.UpdateUserImage(user)
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
	}

	return nil
}
