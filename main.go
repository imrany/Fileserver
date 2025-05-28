package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/imrany/fileserver/config"
)

type Message struct{
	Message string `json:"message"`
}

type DownloadFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Mime string `json:"mime"`
	IsDir bool   `json:"is_dir"`
	FilePath string `json:"file_path"`
	RelativeFilePath string `json:"relative_path"`
}

var downloadPath string

//go:embed templates/*.html
var templates embed.FS

//go:embed static/*
var staticFolder embed.FS

func uploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	fileName := r.FormValue("fileName")
	chunkIndex := r.FormValue("chunkIndex")
	totalChunks := r.FormValue("totalChunks")

	type ErrorResponse struct {
		Error string `json:"error"`
	}
	if fileName == "" {
		errorResponse := ErrorResponse{Error: "File name is required"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		errorResponse := ErrorResponse{Error: "No files received"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	filePath := downloadPath + "/" + fileName
	for _, fh := range files {
		src, err := fh.Open()
		if err != nil {
			http.Error(w, "Error opening uploaded chunk", http.StatusInternalServerError)
			return
		}
		defer src.Close()

		dst, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			http.Error(w, "Unable to write chunk", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = dst.ReadFrom(src)
		if err != nil {
			http.Error(w, "Error saving chunk", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Uploaded chunk %s/%s", chunkIndex, totalChunks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"msg": "Chunk " + chunkIndex + "/" + totalChunks + " received successfully!",
	})
}

func renderIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFS(templates,"templates/*.html") 
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Message string
	}{
		Title:   "File Server",
		Message: "Upload Files",
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "index.html", data)
}

func readDownloads(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.New("").ParseFS(templates,"templates/*.html") 
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	downloadFolder, err := os.ReadDir(downloadPath)
	if err != nil {
		log.Printf("Error accessing downloads: %v", err.Error())
		os.Mkdir(downloadPath, 0755)
		downloadFolder, _ = os.ReadDir(downloadPath)
	}
	downloadFiles := []DownloadFile{}
	for _, file := range downloadFolder {
		fileInfo, err := file.Info()
		if err != nil {
			log.Printf("Error getting file info: %v", err)
			continue
		}

		downloadFiles = append(downloadFiles, DownloadFile{
			Name:  fileInfo.Name(),
			Size:  fileInfo.Size(),
			Mime:  http.DetectContentType([]byte{}), // Placeholder, as we don't have the actual content
			IsDir: file.IsDir(),
			RelativeFilePath: downloadPath + "/" + fileInfo.Name(),
		})
	}

	data := struct {
		Title   string
		Message string
		DownloadFiles []DownloadFile
	}{
		Title:   "Downloads",
		Message: "Download Files",
		DownloadFiles: downloadFiles,
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "downloads.html", data)
}

func main() {
	pwd, err := os.Getwd()
	downloadPath = pwd + "/downloads"
	if err != nil {
		log.Fatalf("Failed to get download folder path %v", err.Error())
	}

	if err := os.Mkdir(downloadPath, 0755); err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create download directory: %v", err)
	}

	staticFs := http.FileServerFS(staticFolder)
	downloadsFs := http.FileServer(http.Dir(downloadPath))

	router := http.NewServeMux()
	router.Handle("GET /static/", staticFs)
	router.Handle("GET /downloads/", http.StripPrefix("/downloads/", downloadsFs))

	router.HandleFunc("GET /", renderIndex)
	router.HandleFunc("GET /downloads", readDownloads)
	router.HandleFunc("POST /api/upload", uploadFile)

	port, err := config.Getenv("PORT")
	if err != nil {
		port = "8080"
	}

	// Open browser after server starts
	go func() {
		url := "http://localhost:" + port
		// Linux, macOS, Windows support
		var cmd string
		var args []string
		switch {
		case os.Getenv("WSL_DISTRO_NAME") != "":
			cmd = "wslview"
			args = []string{url}
		case os.Getenv("XDG_SESSION_TYPE") != "":
			cmd = "xdg-open"
			args = []string{url}
		case os.Getenv("OS") == "Windows_NT":
			cmd = "rundll32"
			args = []string{"url.dll,FileProtocolHandler", url}
		default:
			cmd = "open"
			args = []string{url}
		}
		_ = exec.Command(cmd, args...).Start()
	}()

	srv := http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: router,
	}

	log.Printf("Server running on PORT %v", port)
	srv.ListenAndServe()
}
