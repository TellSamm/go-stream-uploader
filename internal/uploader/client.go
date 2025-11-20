package uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Uploader struct {
	URL     string
	Timeout int
}

func (u *Uploader) UploadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	contentType := mw.FormDataContentType()

	errChan := make(chan error, 1)

	go func() {
		defer pw.Close()

		part, err := mw.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			errChan <- err
			return
		}

		if _, err = io.Copy(part, file); err != nil {
			errChan <- err
			return
		}

		if err = mw.Close(); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	req, err := http.NewRequest("POST", u.URL, pr)
	if err != nil {
		return fmt.Errorf("cannot create request: %v", err)
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{
		Timeout: time.Duration(u.Timeout) * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {

		select {
		case writeErr := <-errChan:
			if writeErr != nil {
				return fmt.Errorf("upload failed (write: %v): %v", writeErr, err)
			}
		default:
		}
		return fmt.Errorf("upload failed: %v", err)
	}
	defer resp.Body.Close()

	if writeErr := <-errChan; writeErr != nil {
		return fmt.Errorf("failed to write multipart body: %v", writeErr)
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Server response:", string(body))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
