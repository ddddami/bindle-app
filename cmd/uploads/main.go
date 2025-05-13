package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ddddami/bindle/uploads"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /files/", http.StripPrefix("/files", http.FileServer(http.Dir("./uploads/"))))

	mux.HandleFunc("GET /{$}", serveUploadForm)
	mux.HandleFunc("POST /upload", handleSingleFileUpload)
	mux.HandleFunc("POST /upload-multiple", handleMultipleFileUpload)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleSingleFileUpload(w http.ResponseWriter, r *http.Request) {
	opts := uploads.FileUploadOptions{
		DestinationDir:    "./uploads",
		MaxSize:           5 * 1024 * 1024, // 5MB
		AllowedExts:       []string{"jpg", "jpeg", "png", "gif", "pdf"},
		RandomizeFilename: true,
		FilenamePrefix:    "upload_",
	}

	err := r.ParseMultipartForm(int64(opts.MaxSize))
	if err != nil {
		fmt.Fprintf(w, "file is too big!")
	}

	savedFile, err := uploads.SaveSingleFormFile(r, "file", &opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upload failed: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h2>File uploaded successfully!</h2>")
	fmt.Fprintf(w, "<p>Original name: %s</p>", savedFile.OriginalName)
	fmt.Fprintf(w, "<p>Saved as: %s</p>", savedFile.SavedName)
	fmt.Fprintf(w, "<p>Size: %d bytes</p>", savedFile.Size)

	ext := filepath.Ext(savedFile.SavedName)
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
		fmt.Fprintf(w, "<p><img src=\"/files/%s\" style=\"max-width:300px;\"></p>", savedFile.SavedName)
	}

	fmt.Fprintf(w, "<p><a href=\"/\">Upload another file</a></p>")
}

func handleMultipleFileUpload(w http.ResponseWriter, r *http.Request) {
	opts := uploads.FileUploadOptions{
		DestinationDir:    "./uploads",
		MaxSize:           5 * 1024 * 1024, // 5MB
		AllowedExts:       []string{"jpg", "jpeg", "png", "gif", "pdf"},
		RandomizeFilename: true,
		FilenamePrefix:    "multi_",
	}

	savedFiles, err := uploads.SaveMultipleFormFiles(r, "files", &opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upload failed: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h2>%d files uploaded successfully!</h2>", len(savedFiles))

	for i, file := range savedFiles {
		fmt.Fprintf(w, "<h3>File %d</h3>", i+1)
		fmt.Fprintf(w, "<p>Original name: %s</p>", file.OriginalName)
		fmt.Fprintf(w, "<p>Saved as: %s</p>", file.SavedName)
		fmt.Fprintf(w, "<p>Size: %d bytes</p>", file.Size)

		// Display the image if it's an image
		ext := filepath.Ext(file.SavedName)
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
			fmt.Fprintf(w, "<p><img src=\"/files/%s\" style=\"max-width:200px;\"></p>", file.SavedName)
		}

		fmt.Fprintf(w, "<hr>")
	}

	fmt.Fprintf(w, "<p><a href=\"/\">Upload more files</a></p>")
}

func serveUploadForm(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>File Upload Example</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .form-container { margin-bottom: 30px; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        h2 { color: #333; }
        input[type="file"] { margin: 10px 0; }
        input[type="submit"] { background: #4CAF50; color: white; padding: 10px 15px; border: none; cursor: pointer; }
    </style>
</head>
<body>
    <h1>File Upload Demo</h1>
    
    <div class="form-container">
        <h2>Single File Upload</h2>
        <form action="/upload" method="post" enctype="multipart/form-data">
            <p>Select a file to upload (max 5MB, only images and PDFs):</p>
            <input type="file" name="file" required>
            <br>
            <input type="submit" value="Upload File">
        </form>
    </div>
    
    <div class="form-container">
        <h2>Multiple File Upload</h2>
        <form action="/upload-multiple" method="post" enctype="multipart/form-data">
            <p>Select multiple files to upload (max 10MB each):</p>
            <input type="file" name="files" multiple required>
            <br>
            <input type="submit" value="Upload Files">
        </form>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
