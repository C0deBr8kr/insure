package crons

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	// "github.com/rk/go-cron"
	"log"
	// "time"
	"bitbucket.org/bdraff/insure/database"
	"github.com/williballenthin/govt"
	"bitbucket.org/bdraff/insure/models"
)

var apikey string
var apiurl string
var file string

// init - initializes flag variables.
func setVar() {
	flag.StringVar(&apikey, "apikey", os.Getenv("VT_API_KEY"), "Set environment variable VT_API_KEY to your VT API Key or specify on prompt")
	flag.StringVar(&apiurl, "apiurl", "https://www.virustotal.com/vtapi/v2/", "URL of the VirusTotal API to be used.")
}

func Init() {
  // cron.NewCronJob(-1, -1, -1, -1, 30, 0, func (time.Time) {
  // 	log.Print("cron")
  // 	setVar()
  // 	retrieve_files_to_be_analyzed()
  // })
  	setVar()
  	analyzeNewBinaries()
  	fetchResultsForSubmittedBinaries()
}

// check - an error checking function
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Send each new binary to be analysed to Virus total and update the 
// permalink in the database.
func sendFileForAnalysis(id int,filename string) {
	if filename == "" {
		fmt.Println("-filename=<fileToScan.ext> missing!")
		os.Exit(1)
	}
	c, err := govt.New(govt.SetApikey(apikey), govt.SetUrl(apiurl))
	check(err)

	// get a filename report
	r, err := c.ScanFile(filename)
	check(err)

	// Marshal into an object
	var res models.VirusScanResponse
	//fmt.Printf("r: %s\n", r)
	j, err := json.Marshal(r)
	check(err)
	
	json.Unmarshal(j, &res)
	check(err)
    
    database.UpdatePermaLink(id, res.Resource, res.PermaLink)
	
	// fmt.Printf("FileReport: ")
	// os.Stdout.Write(j)
}

// Retrieve all new files that are to be analyzed
func analyzeNewBinaries() {
	binariesList := database.GetAllNewBinaries()
	for i := 0; i < len(binariesList); i++ {
		log.Print("Making a call to Virus total for this binary : " + binariesList[i].Path)
		sendFileForAnalysis(binariesList[i].Id, binariesList[i].Path)
    }
}

// Send each new binary to be analysed to Virus total and update the 
// permalink in the database.
func getAnalysisResults(id int, resource string) {
	if resource == "" {
		fmt.Println("-resource=<md5|sha-1|sha-2> not given!")
		os.Exit(1)
	}
	c, err := govt.New(govt.SetApikey(apikey), govt.SetUrl(apiurl))
	check(err)

	// get a file report
	r, err := c.GetFileReport(resource)
	check(err)

	j, err := json.Marshal(r)
	check(err)

	// Marshal into an object
	var res models.VirusScanResponse
	json.Unmarshal(j, &res)
	check(err)

	isSafe := 'N'
	if res.Positives == 0 {
		isSafe = 'Y'
	}

	database.UpdateAnalysisResult(id,isSafe)
	fmt.Printf("FileReport: ")
	os.Stdout.Write(j)
}

// Retrieve all new files that are to be analyzed
func fetchResultsForSubmittedBinaries() {
	binariesList := database.GetAllBinariesAwaitingAnalysisResults()
	for i := 0; i < len(binariesList); i++ {
		log.Print("Making a call to Virus total for this binary to check if analysis is done: " + binariesList[i].Path)
		getAnalysisResults(binariesList[i].Id, binariesList[i].Path)
    }
}