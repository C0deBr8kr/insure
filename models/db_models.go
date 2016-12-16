package models

import ("time")
// Not exported so that we can also provide default values to initialised objects
type UploadedFile struct {
	Id int
	Name string
	Email string
	Path string
	Time time.Time
	IsVirusChecked bool
}

// Acts as a constructor uploadedFile struct
func NewUploadedFile() UploadedFile {
	newUploaded := UploadedFile{}
	return newUploaded
}