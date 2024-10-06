package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type FileUpload struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Path string    `json:"path"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./static/index.html")
}

func UploadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Parse the Input, Type multipart/form-data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//Retrive file from form data

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from form-data: "+err.Error(), http.StatusInternalServerError)
		fmt.Printf("Failed to retrieve file from form-data: %s", err.Error())
		return
	}
	defer file.Close()

	// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v", handler.Size)
	// fmt.Printf("File Header: %+v", handler.Header)

	//Create a file and write the contents of uploaded file to it

	extention := filepath.Ext(handler.Filename)

	fmt.Println("file extention: " + extention)

	fileName := strings.TrimSuffix(handler.Filename, extention)

	newFileFullName := fmt.Sprintf("%s_*%s", fileName, extention)

	fmt.Println("New file FullName: " + newFileFullName)

	newFile, err := os.CreateTemp("uploads", newFileFullName)
	if err != nil {
		http.Error(w, "Failed to create a new file on server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	uploadedFileContentInBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to create a new file on server: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newFile.Write(uploadedFileContentInBytes)

	// Sample data to be written to JSON
	filePath := newFile.Name()
	newFileUpload := FileUpload{
		Id:   uuid.New(),
		Name: strings.Split(filePath, "/")[1],
		Path: filePath,
	}

	// Read the existing JSON file
	jsonFileName := "file.json"
	var files []FileUpload

	if _, err := os.Stat(jsonFileName); err == nil {
		// File exists, read it

		jsonFileData, err := os.ReadFile(jsonFileName)
		if err != nil {
			http.Error(w, "Error reading json file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the existing data into the Go slice
		err = json.Unmarshal(jsonFileData, &files)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}

	}

	// Append the new data
	files = append(files, newFileUpload)

	// Convert the Go data structure to JSON
	updatedJsonData, err := json.MarshalIndent(files, "", " ")
	if err != nil {
		http.Error(w, "Error marshalling  updated json data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the updated JSON data back to the file
	err = os.WriteFile(jsonFileName, updatedJsonData, 0644)
	if err != nil {
		http.Error(w, "Error writing updated json to file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "File Uploaded Successfully \n")
}

func DownloadsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonFileName := "file.json"
	var files []FileUpload

	w.Header().Set("Content-Type", "application/json")

	if _, err := os.Stat(jsonFileName); err == nil {
		// File exists, read it

		jsonFileData, err := os.ReadFile(jsonFileName)
		if err != nil {
			http.Error(w, "Error reading json file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the existing data into the Go slice
		err = json.Unmarshal(jsonFileData, &files)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(files)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	parsedUUID, err := uuid.Parse(ps.ByName("fileId"))
	if err != nil {
		http.Error(w, "Error parsing UUID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonFileName := "file.json"
	var files []FileUpload

	w.Header().Set("Content-Type", "application/json")

	if _, err := os.Stat(jsonFileName); err == nil {
		// File exists, read it

		jsonFileData, err := os.ReadFile(jsonFileName)
		if err != nil {
			http.Error(w, "Error reading json file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the existing data into the Go slice
		err = json.Unmarshal(jsonFileData, &files)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var requestedFile FileUpload

	for _, file := range files {
		if file.Id == parsedUUID {
			requestedFile = file
			break
		}
	}
	json.NewEncoder(w).Encode(requestedFile)
}

func SetupRoutes() {
	serverPort := 8080
	addr := fmt.Sprintf(":%d", serverPort)

	// Create a new router
	router := httprouter.New()

	// Define routes
	router.GET("/", HomeHandler)
	router.POST("/upload", UploadHandler)
	router.GET("/downloads", DownloadsHandler)
	router.GET("/downloads/:fileId", DownloadHandler)

	// Start the server
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}

}

func main() {
	SetupRoutes()
}
