package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"bitbucket.org/bdraff/insure/database"
	"bitbucket.org/bdraff/insure/utils"
	"bitbucket.org/bdraff/insure/models"
	"bitbucket.org/bdraff/insure/crons"
	"log"
)

func init() {
	// ALso have the line where the log was made.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	crons.Init()
}

// Uploads the binary to the uploaded_files folder
func uploadFile(w http.ResponseWriter, r *http.Request) models.UploadedFile{
	// log.Print(r.ContentType)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		// fmt.Fprintln(w, err)
		log.Print(err)
		// os.Exit(1)
	}
	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")

	// check if the request is a POST request or not
	if err != nil {
		// fmt.Fprintln(w, err)
		log.Print(err)
		// os.Exit(1)
	}

	file.Close()

	filename := utils.GenerateSavedFilename(header.Filename)
	fullyQualifiedName := "uploaded_files/"+ filename
	//creates a new file with name "uploadedfile"
	out, err := os.Create(fullyQualifiedName)
	if err != nil {
		log.Print("File could not be created")
		log.Print(err)
		// os.Exit(1)
	}

	defer out.Close()
	// write the content from POST to the file
	_, err = io.Copy(out, file)
	
	if err != nil {
		log.Print(err)
		// fmt.Fprintln(w, err)
		// log.Fatal(err)
	}

	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, filename)

	newUploadedFile := models.NewUploadedFile()
	newUploadedFile.Name = filename
	newUploadedFile.Path = fullyQualifiedName
	return newUploadedFile
}

// Extract data from the posted form
func getFormData(w http.ResponseWriter, r *http.Request, newUploadedFile *models.UploadedFile) {
	r.ParseForm()
	newUploadedFile.Email = r.Form["email"][0]
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	newUploadedFile := uploadFile(w, r)
	getFormData(w, r, &newUploadedFile)
	database.InsertNewBinary(&newUploadedFile)
}

func main() {
	http.HandleFunc("/", handlePostRequest)
	http.ListenAndServe(":8080", nil)
	
}
