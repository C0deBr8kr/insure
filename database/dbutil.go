package database

import (
	"bitbucket.org/bdraff/insure/config"
	"bitbucket.org/bdraff/insure/models"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
)

var db *sql.DB

func connect () {
	var configuration = config.GetConfig()
	fmt.Println(configuration.Username + ":" + configuration.Password)
	var err error
	db, err = sql.Open("mysql",
		configuration.Username+":"+configuration.Password+"@tcp(127.0.0.1:3306)/insure")
	if err != nil {
		log.Fatal(err)
	}

	// defer db.Close()
	db.SetMaxIdleConns(100)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connected")
}

func disconnect() {
	db.Close()
}

func InsertNewBinary(uploadedFile *models.UploadedFile) {
	connect()
	log.Print("inserted into database")
	log.Print(uploadedFile.Name + ":" + uploadedFile.Email)
	stmt, err := db.Prepare("INSERT INTO binaries(name,path,email) VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(uploadedFile.Name, uploadedFile.Path, uploadedFile.Email)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	disconnect()
}

func GetAllNewBinaries() []models.UploadedFile {
	connect()
	rows, err := db.Query("select id, path from binaries where is_virus_scanned = 'N' and vt_permalink is null")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	binaryList := []models.UploadedFile{}
    for rows.Next() {
            var r models.UploadedFile
            err = rows.Scan(&r.Id, &r.Path)
            if err != nil {
                    log.Fatalf("Scan: %v", err)
            }
            binaryList = append(binaryList, r)
    }
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	disconnect()
	return binaryList
}

func UpdatePermaLink(id int, resource string, permalink string) {
	connect()
	stmt, err := db.Prepare("UPDATE binaries set vt_permalink=?, vt_resourceid=? where id=?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(permalink, resource, id)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	disconnect()
}

func GetAllBinariesAwaitingAnalysisResults() []models.UploadedFile {
	connect()
	rows, err := db.Query("select id, vt_resourceid from binaries where is_virus_scanned = 'N' and vt_permalink is not null")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	binaryList := []models.UploadedFile{}
    for rows.Next() {
            var r models.UploadedFile
            err = rows.Scan(&r.Id, &r.Path)
            if err != nil {
                    log.Fatalf("Scan: %v", err)
            }
            binaryList = append(binaryList, r)
    }
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	disconnect()
	return binaryList
}

func UpdateAnalysisResult(id int, isSafe byte) {
	connect()
	stmt, err := db.Prepare("UPDATE binaries set is_safe=?,is_virus_scanned='Y' where id=?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(isSafe, id)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	disconnect()
}