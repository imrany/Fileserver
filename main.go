package main

import (
	"net/http"
	"os"
	"log"
	"encoding/json"
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
}

var downloadPath string

func helloWorldJson(w http.ResponseWriter , r *http.Request){
	helloMsg :=Message{
		Message:"Hello world",
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(helloMsg)
}

func uploadFile(w http.ResponseWriter, r *http.Request){
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create(downloadPath + "/" + handler.Filename)
	if err != nil {
		http.Error(w, "Unable to create the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = dst.ReadFrom(file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Message{Message: "File uploaded successfully"})
}

func readDownloads(w http.ResponseWriter, r *http.Request){
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

		pwd, err := os.Getwd()
		filePath := pwd + downloadPath + "/" + fileInfo.Name()
		if err != nil {
			filePath = ""
		}
		downloadFiles = append(downloadFiles, DownloadFile{
			Name:  fileInfo.Name(),
			Size:  fileInfo.Size(),
			Mime:  http.DetectContentType([]byte{}), // Placeholder, as we don't have the actual content
			IsDir: file.IsDir(),
			FilePath: filePath,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(downloadFiles)
}

func main(){
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	downloadPath =pwd + "/downloads"

	if err := os.Mkdir(downloadPath, 0755); err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create download directory: %v", err)
	}


	fs := http.FileServer(http.Dir("./views"))

	router := http.NewServeMux()
	router.Handle("GET /views/", http.StripPrefix("/views/", fs))
	router.HandleFunc("GET /", helloWorldJson)
	router.HandleFunc("POST /upload", uploadFile)
	router.HandleFunc("GET /read_file", readDownloads)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := http.Server{
		Addr: "0.0.0.0:" + port,
		Handler: router,
	}

	log.Printf("Server running on PORT %v", port)
	srv.ListenAndServe()
}
