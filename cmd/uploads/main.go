package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ddddami/bindle/uploads"
)

func main() {
	// Create a test file for download demo
	uploadDir := "./uploads"
	testFilePath := filepath.Join(uploadDir, "test-document.pdf")
	createTestFile(testFilePath)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home)

	mux.Handle("GET /files/", http.StripPrefix("/files", http.FileServer(http.Dir("./uploads/"))))

	mux.HandleFunc("GET /uploads", serveUploadForm)
	mux.HandleFunc("POST /upload", handleSingleFileUpload)
	mux.HandleFunc("POST /upload-multiple", handleMultipleFileUpload)

	mux.HandleFunc("GET /downloads", serveDownloadPage)
	mux.HandleFunc("GET /download", handleSingleDownload)
	mux.HandleFunc("GET /download/with-name", handleDownloadWithCustomName)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func home(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Bindle Demo Usage</title>

  <style>
  body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; font-size: 1.2rem;}
  </style>
</head>
<body>
    <h1>File Upload Demo</h1>
    
  <ul>
  <li>Downloads <a href="/downloads">downloads page</a></li>
  <li>Uploads <a href="/uploads">uploads page</a></li>
  </ul>
    
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
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

func handleSingleDownload(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("./uploads", "wallpaper_6.jpg")

	if err := uploads.ServeFileForDownload(w, r, filePath, uploads.DownloadOptions{ForceDownload: true, ExtraHeaders: make(map[string]string)}); err != nil {
		http.Error(w, "Error serving file: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleDownloadWithCustomName(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("./uploads", "test-document.pdf")

	opts := uploads.DownloadOptions{
		ForceDownload:     true,
		SuggestedFilename: "renamed-document.pdf",
		ExtraHeaders: map[string]string{
			"X-Download-Type": "Custom",
		},
	}

	if err := uploads.ServeFileForDownload(w, r, filePath, opts); err != nil {
		http.Error(w, "Error serving file: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func serveDownloadPage(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>File Download Examples</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .section { margin-bottom: 30px; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        h2 { color: #333; }
        ul { list-style-type: none; padding: 0; }
        li { margin: 10px 0; }
        a { color: #0066cc; text-decoration: none; padding: 8px 16px; background-color: #f0f0f0; 
            border-radius: 4px; display: inline-block; }
        a:hover { background-color: #e0e0e0; }
        .note { font-size: 0.9em; color: #666; margin-top: 5px; }
    </style>
</head>
<body>
    <h1>File Download Examples</h1>
    
    <div class="section">
        <h2>Download Links</h2>
        <ul>
            <li><a href="/download">Download Test Document</a>
                <div class="note">Uses default download options</div>
            </li>
            <li><a href="/download/with-name">Download with Custom Filename</a>
                <div class="note">Changes the suggested filename to "renamed-document.pdf"</div>
            </li>
        </ul>
    </div>
    
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func createTestFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	// Write some dummy content that looks like a PDF header
	file.WriteString("%PDF-1.4\n")
	file.WriteString("1 0 obj\n")
	file.WriteString("<< /Type /Catalog /Pages 2 0 R >>\n")
	file.WriteString("endobj\n")
	file.WriteString("This is a test document for download functionality.\n")
	file.WriteString("%%EOF\n")
}
