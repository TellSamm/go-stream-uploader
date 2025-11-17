package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Run() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		out, err := os.Create("uploaded_" + header.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Upload OK: %s\n", header.Filename)
	})

	fmt.Println("Server started on http://localhost:8080/upload")
	http.ListenAndServe(":8080", nil)
}
