package uploader

func NewUploader(url string) *Uploader {
	return &Uploader{
		URL:     url,
		Timeout: 30,
	}
}
