package utils

import (
	"time"
	"strings"
	"fmt"
)

type file struct {
	filename string
	extension string
}

// Seperate the filename and the extension
func getFilename(complete_filename string) file{
	file_obj := file{}
	if(strings.Contains(complete_filename,".")) {
		file_obj.filename = complete_filename[:strings.LastIndex(complete_filename, ".")]
		file_obj.extension = complete_filename[strings.LastIndex(complete_filename, ".")+1:]
	} else {
		file_obj.filename = complete_filename
	}
	return file_obj
}

func GenerateSavedFilename(complete_filename string) string{
	file_obj := getFilename(complete_filename)
	newFilename := file_obj.filename + "_" + time.Now().UTC().Format("2006010215040") + 
			"." + file_obj.extension
	fmt.Println(newFilename)
	return newFilename
}