package uploader

import "time"

func NewUploader(url string) *Uploader {
	return &Uploader{
		URL:            url,
		ConnectTimeout: 30 * time.Second,
		UploadTimeout:  0, //помечу (тело запроса может литься сколько угодно ему)
	}
}
