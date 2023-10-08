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
	handleUploadUserImages(w http.ResponseWriter, r *http.Request) error
}

func (s *ApiRouter) handleUploadUserImages(w http.ResponseWriter, r *http.Request) {

	id, err := getID(r)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	user, err := s.store.GetUser(id)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	if err := r.ParseMultipartForm(32 * 1024 * 1024); err != nil {
		s.handleError(w, r, err)
		return
	}

	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			s.handleError(w, r, err)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" {
			s.handleError(w, r, err)
			return
		}

		err = os.MkdirAll("./static/uploads", os.ModePerm)
		if err != nil {
			s.handleError(w, r, err)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			s.handleError(w, r, err)
			return
		}

		newFileName := fmt.Sprintf("./static/uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		f, err := os.Create(newFileName)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		if len(user.ImageURL) > 2 {
			err := os.Remove(user.ImageURL[1:])
			if err != nil {
				s.handleError(w, r, err)
				return
			}
		}
		user.ImageURL = "." + newFileName
		s.store.UpdateUserImage(user)

		defer f.Close()
		user, err := s.store.GetUser(id)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		_, err = io.Copy(f, file)
		if err != nil {
			s.handleError(w, r, err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte(user.ImageURL))
		if err != nil {
			s.handleError(w, r, err)
			return
		}
	}
}
