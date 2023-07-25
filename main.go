package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func pdftotextHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the PDF file content.
	file, _, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Failed to parse PDF file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary file to save the uploaded PDF.
	tmpFile, err := ioutil.TempFile("", "input_*.pdf")
	if err != nil {
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Save the PDF content to the temporary file.
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, "Failed to save PDF content to temporary file", http.StatusInternalServerError)
		return
	}

	// Execute pdftotext command.
	cmd := exec.Command("pdftotext", tmpFile.Name(), "-")
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to convert PDF to text", http.StatusInternalServerError)
		return
	}

	// Set response content type to plain text.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Write the PDF text to the response.
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Register the /pdftotxt endpoint to the pdftotextHandler function.
	http.HandleFunc("/pdftotxt", pdftotextHandler)

	// Start the server on port 8080.
	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
