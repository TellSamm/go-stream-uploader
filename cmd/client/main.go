package client

import (
	"fmt"
	"go-stream-uploader/internal/uploader"
)

func main() {
	u := uploader.NewUploader("http://localhost:8080/upload")

	err := u.UploadFile("testfiles/bigfile.bin")
	if err != nil {
		fmt.Println("Upload error:", err)
	}
}
