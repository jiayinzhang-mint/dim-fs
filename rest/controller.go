package rest

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Fetch image file
func viewImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	path := params["path"]

	fullPath := viper.GetString("file.upload") + path
	fileinfo, e := os.Stat(fullPath)

	// File not exist
	if os.IsNotExist(e) {
		logrus.Error(path, " does not exists.")
		w.WriteHeader(404)
		return
	}

	// Path is a dir
	if fileinfo.IsDir() {
		w.WriteHeader(400)
		return
	}

	// Check if file exists and open
	Openfile, err := os.Open(fullPath)
	defer Openfile.Close() // Close after function return
	if err != nil {
		// File not found, send 404
		http.Error(w, "File not found.", 404)
	}
	http.ServeFile(w, r, fullPath)
}

// Upload image file
func uploadImage(w http.ResponseWriter, r *http.Request) {

	// Upload path
	targetPath := r.FormValue("targetPath")
	uploadDir := viper.GetString("file.upload")

	// File
	uploadFile, header, getFormErr := r.FormFile("image")
	if getFormErr != nil {
		w.WriteHeader(400)
		return
	}

	// Create path if not exist
	if _, notExist := os.Stat(filepath.Join(uploadDir, targetPath)); os.IsNotExist(notExist) {
		os.MkdirAll(filepath.Join(uploadDir, targetPath), 0777)
	}

	// Create file
	saveFile, createErr := os.OpenFile(filepath.Join(uploadDir, targetPath, header.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if createErr != nil {
		logrus.Error(createErr, "Failed to create file")
		w.WriteHeader(500)
		return
	}

	// Write file
	_, wtErr := io.Copy(saveFile, uploadFile)
	if wtErr != nil {
		logrus.Error(wtErr, "Failed to write file")
		w.WriteHeader(500)
		return
	}

	defer saveFile.Close()
	defer uploadFile.Close()

	w.WriteHeader(200)
	return
}
