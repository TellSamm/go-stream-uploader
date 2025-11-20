package uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Uploader struct {
	URL     string
	Timeout int
}

func (u *Uploader) UploadFile(filePath string) error {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	contentType := mw.FormDataContentType()

	go func() {
		defer pw.Close()

		part, err := mw.CreateFormFile("file", contentType)
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}

		if _, err = io.Copy(part, file); err != nil {
			_ = pw.CloseWithError(err)
			return
		}

		if err := mw.Close(); err != nil {
			_ = pw.CloseWithError(err)
			return
		}
	}()

	req, err := http.NewRequest("POST", u.URL, pr)
	if err != nil {
		return fmt.Errorf("cannot create request: %v", err)
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())

	client := &http.Client{Timeout: time.Duration(u.Timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Server response:", string(body))

	return nil
}
