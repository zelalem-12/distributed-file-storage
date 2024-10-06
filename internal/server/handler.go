package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// "myproject/internal/services"
	"context"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/zelalem-12/distributed-file-storage/internal/domain"
	"github.com/zelalem-12/distributed-file-storage/internal/services"
)

type FileHandler struct {
	FileService services.FileServiceInterface
}

func NewFileHandler(fileService services.FileServiceInterface) *FileHandler {
	return &FileHandler{
		FileService: fileService,
	}
}

func (h *FileHandler) HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./static/index.html")
}

func (h *FileHandler) GetFileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	parsedUUID, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Error parsing UUID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := h.FileService.GetById(context.Background(), parsedUUID)
	if err != nil {
		http.Error(w, "Error parsing UUID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	responseDTO := &Response{
		Id:        file.GetID(),
		Name:      file.GetName(),
		Path:      file.GetPath(),
		CreatedAt: file.GetCreatedAt().Format("006-01-02 15:04:05"),
	}
	json.NewEncoder(w).Encode(responseDTO)
}

func (h *FileHandler) UploadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Parse the Input, Type multipart/form-data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//Retrive file from form data

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from form-data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	extention := filepath.Ext(handler.Filename)

	fileName := strings.TrimSuffix(handler.Filename, extention)

	newFileFullName := fmt.Sprintf("%s_*%s", fileName, extention)

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

	filePath := newFile.Name()
	newFileObject, err := domain.CreateFile(strings.Split(filePath, "/")[1], filePath)

	if err != nil {
		http.Error(w, "Failed to create a new file object: "+err.Error(), http.StatusInternalServerError)
		return
	}

	savedData, err := h.FileService.Create(context.Background(), newFileObject)
	if err != nil {
		http.Error(w, "Failed to persist a new file : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Sample data to be written to JSON

	// newFileUpload := FileUpload{
	// 	Id:        uuid.New(),
	// 	Name:      strings.Split(filePath, "/")[1],
	// 	Path:      filePath,
	// 	CreatedAt: time.Now(),
	// }

	// Read the existing JSON file
	// jsonFileName := "file.json"
	// var files []FileUpload

	// if _, err := os.Stat(jsonFileName); err == nil {
	// 	// File exists, read it

	// 	jsonFileData, err := os.ReadFile(jsonFileName)
	// 	if err != nil {
	// 		http.Error(w, "Error reading json file: "+err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Unmarshal the existing data into the Go slice
	// 	err = json.Unmarshal(jsonFileData, &files)
	// 	if err != nil {
	// 		http.Error(w, "Error unmarshaling JSON: "+err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// }

	// // Append the new data
	// files = append(files, newFileUpload)

	// // Convert the Go data structure to JSON
	// updatedJsonData, err := json.MarshalIndent(files, "", " ")
	// if err != nil {
	// 	http.Error(w, "Error marshalling  updated json data: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// // Write the updated JSON data back to the file
	// err = os.WriteFile(jsonFileName, updatedJsonData, 0644)
	// if err != nil {
	// 	http.Error(w, "Error writing updated json to file: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	fmt.Fprint(w, savedData.GetID())
}
func (h *FileHandler) DownloadsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//Read from Json file and return back to the client

	//jsonFileName := "file.json"
	//var files []FileUpload

	// if _, err := os.Stat(jsonFileName); err == nil {
	// 	// File exists, read it

	// 	jsonFileData, err := os.ReadFile(jsonFileName)
	// 	if err != nil {
	// 		http.Error(w, "Error reading json file: "+err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Unmarshal the existing data into the Go slice
	// 	err = json.Unmarshal(jsonFileData, &files)
	// 	if err != nil {
	// 		http.Error(w, "Error unmarshaling JSON: "+err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	//}
	files, err := h.FileService.GetAll(context.Background())
	if err != nil {
		http.Error(w, "Error parsing UUID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	response := []Response{}
	for _, file := range files {
		response = append(response, Response{
			Id:        file.GetID(),
			Name:      file.GetName(),
			Path:      file.GetPath(),
			CreatedAt: file.GetCreatedAt().Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
