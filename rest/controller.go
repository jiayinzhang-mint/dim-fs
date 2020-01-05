package rest

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func viewImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	path := params["path"]

	fullPath := viper.GetString("file.upload") + path
	fileinfo, e := os.Stat(fullPath)

	// File not exist
	if os.IsNotExist(e) {
		log.Println(path, " does not exists.")
		w.WriteHeader(404)
		return
	}

	// Path is a dir
	if fileinfo.IsDir() {
		w.WriteHeader(500)
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
