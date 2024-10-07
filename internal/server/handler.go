package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/zelalem-12/distributed-file-storage/internal/domain"
	"github.com/zelalem-12/distributed-file-storage/internal/services"
	"github.com/zelalem-12/distributed-file-storage/internal/utils"
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

// UploadFilesHandler godoc
// @Summary Upload multiple files
// @Description Upload multiple files to the distributed file storage system
// @Tags files
// @Accept  mpfd
// @Produce json
// @Param files formData file true "Files to upload"
// @Success 200 {array} string "List of uploaded file IDs"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /upload [post]
func (h *FileHandler) UploadFilesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if err := r.ParseMultipartForm(1 << 30); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	formFiles := r.MultipartForm.File["files"]
	if len(formFiles) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	fileIDsChan := make(chan uuid.UUID, len(formFiles))
	errorChan := make(chan error, 1)

	var hasErrorOccurred bool
	var errOccurred error

	for _, formFile := range formFiles {
		wg.Add(1)

		go func(formFile *multipart.FileHeader) {
			defer wg.Done()

			if formFile == nil {
				errorChan <- fmt.Errorf("nil formFile detected")
				return
			}

			file, err := formFile.Open()
			if err != nil {
				errorChan <- fmt.Errorf("failed to retrieve file from form-data: %w", err)
				return
			}
			defer file.Close()

			extension := filepath.Ext(formFile.Filename)

			newFile, err := os.CreateTemp("uploads", fmt.Sprintf("%s_*%s", strings.TrimSuffix(formFile.Filename, extension), extension))
			if err != nil {
				errorChan <- fmt.Errorf("failed to create a new file on server: %w", err)
				return
			}
			defer newFile.Close()

			uploadedFileContentInBytes, err := io.ReadAll(file)
			if err != nil {
				errorChan <- fmt.Errorf("failed to read file content: %w", err)
				return
			}

			_, err = newFile.Write(uploadedFileContentInBytes)
			if err != nil {
				errorChan <- fmt.Errorf("failed to write file content to new file: %w", err)
				return
			}
			filePath := newFile.Name()

			fileSize, fileType, fileExtension, err := utils.GetFileInfo(newFile.Name())
			if err != nil {
				errorChan <- fmt.Errorf("failed to get file infomration: %w", err)
				return
			}

			newFileObject, err := domain.CreateFile(strings.Split(filePath, "/")[1], filePath, fileType, fileExtension, fileSize)
			if err != nil {
				errorChan <- fmt.Errorf("failed to create a new file object: %w", err)
				return
			}

			if newFileObject == nil {
				errorChan <- fmt.Errorf("newFileObject is nil")
				return
			}

			savedData, err := h.FileService.Create(context.Background(), newFileObject)
			if err != nil {
				errorChan <- fmt.Errorf("failed to persist a new file: %w", err)
				return
			}

			if savedData == nil {
				errorChan <- fmt.Errorf("savedData is nil")
				return
			}

			fileIDsChan <- savedData.GetID()
		}(formFile)
	}

	go func() {
		wg.Wait()
		close(fileIDsChan)
		close(errorChan)
	}()

	var fileIDs []uuid.UUID

	for fileID := range fileIDsChan {
		fileIDs = append(fileIDs, fileID)
	}

	select {
	case err := <-errorChan:
		if err != nil {
			hasErrorOccurred = true
			errOccurred = err
		}
	default:
	}

	if hasErrorOccurred {
		http.Error(w, errOccurred.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileIDs)
}

// GetFilesDataHandler godoc
// @Summary Get all file metadata
// @Description Retrieve metadata of all uploaded files
// @Tags files
// @Produce json
// @Success 200 {array} Response "List of files"
// @Failure 404 {string} string "No files found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /files [get]
func (h *FileHandler) GetFilesDataHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	files, err := h.FileService.GetAll(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(files) == 0 {
		http.Error(w, "record not found", http.StatusNotFound)
		return
	}
	response := []Response{}
	for _, file := range files {
		response = append(response, Response{
			Id:        file.GetID(),
			Name:      file.GetName(),
			Path:      file.GetPath(),
			Type:      file.GetType(),
			Size:      file.GetSize(),
			Extension: file.GetExtension(),
			CreatedAt: file.GetCreatedAt().Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DownloadFileByIDHandler godoc
// @Summary Download a file by ID
// @Description Download a file from the distributed storage system by its unique ID
// @Tags files
// @Produce octet-stream
// @Param id path string true "File ID"
// @Success 200 {file} file "Downloaded file"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /downloads/{id} [get]
func (h *FileHandler) DownloadFileByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	parsedUUID, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileData, err := h.FileService.GetById(context.Background(), parsedUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(fileData.GetPath())))

	mergedWriter := w

	err = utils.DownloadFileInParallel(fileData, mergedWriter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
