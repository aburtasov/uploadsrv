package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

// Compile templates on start of the application

var templates = template.Must(template.ParseFiles("public/upload.html"))

// Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFiles(w, r)
	}
}

func main() {
	// Upload route
	http.HandleFunc("/upload", uploadHandler)

	//Listen on port 8080
	http.ListenAndServe(":80", nil)
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {

	err := os.Chdir("files")
	if err != nil {
		log.Fatal(err)
	}

	err = r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["myFile"] // grab the filenames

	for i, _ := range files { // loop through the files one by one

		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		out, err := os.Create(files[i].Filename)

		defer out.Close()
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
			return
		}

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		fmt.Fprintf(w, "Files uploaded successfully : ")
		fmt.Fprintf(w, files[i].Filename+"\n")

	}

	err = os.Chdir("../")
	if err != nil {
		log.Fatal(err)
	}
}
