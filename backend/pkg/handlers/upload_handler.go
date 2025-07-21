package handlers

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/google/uuid"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse multipart form
    err := r.ParseMultipartForm(10 << 20) // 10 MB limit
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    file, header, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Unable to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Validate file type
    contentType := header.Header.Get("Content-Type")
    if !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "video/") {
        http.Error(w, "Invalid file type", http.StatusBadRequest)
        return
    }

    // Create uploads directory if it doesn't exist
    uploadsDir := "uploads"
    if err := os.MkdirAll(uploadsDir, 0755); err != nil {
        http.Error(w, "Unable to create uploads directory", http.StatusInternalServerError)
        return
    }

    // Generate unique filename
    ext := filepath.Ext(header.Filename)
    filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
    filepath := filepath.Join(uploadsDir, filename)

    // Create the file
    dst, err := os.Create(filepath)
    if err != nil {
        http.Error(w, "Unable to create file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    // Copy file content
    _, err = io.Copy(dst, file)
    if err != nil {
        http.Error(w, "Unable to save file", http.StatusInternalServerError)
        return
    }

    // Return file URL
    fileURL := fmt.Sprintf("/uploads/%s", filename)
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"url": "%s"}`, fileURL)
}