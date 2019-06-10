package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Filelog string

var DirLog = "log_server_grls"
var DirTemp = "temp_server_grls"
var FileLog Filelog
var Port string
var GrlsTable = "grls"
var GrlsExceptTable = "grls_except"
var SecretKey = "11111111" //change it
var mutex sync.Mutex
var FileDB = "grls.db"
var ArZir = "file.zip"

func DbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=OFF&_synchronous=OFF", FileDB))
	return db, err
}

func GetPort() {
	flag.Parse()
	ArgS := flag.Arg(0)
	if ArgS == "" {
		fmt.Println("server run on default port 8181")
		Port = "8181"
	} else {
		Port = ArgS
	}
}
func Logging(args ...interface{}) {
	mutex.Lock()
	file, err := os.OpenFile(string(FileLog), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println("Ошибка записи в файл лога", err)
		return
	}
	fmt.Fprintf(file, "%v  ", time.Now())
	for _, v := range args {

		fmt.Fprintf(file, " %v", v)
	}
	//fmt.Fprintf(file, " %s", UrlXml)
	fmt.Fprintln(file, "")
	mutex.Unlock()
}
func CreateLogFile() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirlog := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, DirLog))
	if _, err := os.Stat(dirlog); os.IsNotExist(err) {
		err := os.MkdirAll(dirlog, 0711)

		if err != nil {
			fmt.Println("Не могу создать папку для лога")
			os.Exit(1)
		}
	}
	t := time.Now()
	ft := t.Format("2006-01-02")
	FileLog = Filelog(filepath.FromSlash(fmt.Sprintf("%s/log_grls_server_%v.log", dirlog, ft)))
}

func CreateTempDir() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirtemp := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, DirTemp))
	if _, err := os.Stat(dirtemp); os.IsNotExist(err) {
		err := os.MkdirAll(dirtemp, 0711)

		if err != nil {
			fmt.Println("Не могу создать папку для временных файлов")
			os.Exit(1)
		}
	} else {
		err = os.RemoveAll(dirtemp)
		if err != nil {
			fmt.Println("Не могу удалить папку для временных файлов")
		}
		err := os.MkdirAll(dirtemp, 0711)
		if err != nil {
			fmt.Println("Не могу создать папку для временных файлов")
			os.Exit(1)
		}
	}
}
func GetPass() {
	flag.Parse()
	ArgS := strings.TrimSpace(flag.Arg(1))
	if ArgS != "" {
		SecretKey = ArgS
		fmt.Println(fmt.Sprintf("Установлен пароль для обновления базы - %s", SecretKey))
	} else {
		fmt.Println(fmt.Sprintf("Пароль для обновления базы оставлен по умолчанию - %s", SecretKey))
	}
}

func CreateDB() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDB := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, FileDB))
	if _, err := os.Stat(fileDB); os.IsNotExist(err) {
		fmt.Println(err)
		f, err := os.Create(fileDB)
		if err != nil {
			Logging(err)
			panic(err)
		}
		err = f.Chmod(0777)
		if err != nil {
			Logging(err)
			//panic(err)
		}
		err = f.Close()
		if err != nil {
			Logging(err)
			panic(err)
		}
		db, err := DbConnection()
		if err != nil {
			Logging(err)
			panic(err)
		}
		defer db.Close()
		_, err = db.Exec(`CREATE TABLE "grls" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"mnn"	TEXT,
	"name"	TEXT,
	"form"	TEXT,
	"owner"	TEXT,
	"atx"	TEXT,
	"quantity"	INTEGER,
	"max_price"	REAL,
	"first_price"	REAL,
	"ru"	TEXT,
	"date_reg"	TEXT,
	"code"	TEXT,
	"date_pub"	TEXT
)`)
		if err != nil {
			Logging(err)
			panic(err)
		}
		_, err = db.Exec(`CREATE TABLE "grls_except" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"mnn"	TEXT,
	"name"	TEXT,
	"form"	TEXT,
	"owner"	TEXT,
	"atx"	TEXT,
	"quantity"	INTEGER,
	"max_price"	REAL,
	"first_price"	REAL,
	"ru"	TEXT,
	"date_reg"	TEXT,
	"code"	TEXT,
	"except_cause"	TEXT,
	"except_date"	TEXT,
	"date_pub"	TEXT
)`)
		if err != nil {
			Logging(err)
			panic(err)
		}
	}
}

func CreateEnv() {
	CreateLogFile()
	CreateTempDir()
	GetPort()
	GetPass()
	CreateDB()
}
