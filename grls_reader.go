package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/extrame/xls"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
)

type GrlsReader struct {
	Url   string
	Added int
}

func (t *GrlsReader) reader() {
	p := t.downloadString()
	if p == "" {
		Logging("get empty string", p)
		return
	}
	url := t.extractUrl(p)
	if url == "" {
		Logging("get empty url", p)
		return
	}
	t.downloadArchive(url)

}

func (t *GrlsReader) downloadString() string {
	pageSource := DownloadPage(t.Url)
	return pageSource

}

func (t *GrlsReader) extractUrl(p string) string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(p))
	if err != nil {
		Logging(err)
		return ""
	}
	aTag := doc.Find("#ctl00_plate_tdzip > a").First()
	if aTag == nil {
		Logging("a tag not found")
		return ""
	}
	href, ok := aTag.Attr("href")
	if !ok {
		Logging("href attr in a tag not found")
		return ""
	}
	return fmt.Sprintf("https://grls.rosminzdrav.ru/%s", href)
}

func (t *GrlsReader) downloadArchive(url string) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filePath := filepath.FromSlash(fmt.Sprintf("%s/%s/%s", dir, DirTemp, ArZir))
	err := DownloadFile(filePath, url)
	if err != nil {
		Logging("file was not downloaded, exit", err)
		return
	}
	dirZip := filepath.FromSlash(fmt.Sprintf("%s/%s/", dir, DirTemp))
	err = Unzip(filePath, dirZip)
	if err != nil {
		Logging("file was not unzipped, exit", err)
		return
	}
	files, err := FilePathWalkDir(dirZip)
	if err != nil {
		Logging("filelist return error, exit", err)
		return
	}
	for _, f := range files {
		if strings.HasSuffix(f, "xls") {
			t.extractXlsData(f)
		}
	}
}

func (t *GrlsReader) extractXlsData(nameFile string) {
	defer SaveStack()
	xlFile, err := xls.Open(nameFile, "utf-8")
	if err != nil {
		Logging("error open excel file, exit", err)
		return
	}
	sheet := xlFile.GetSheet(0)
	t.insertToBase(sheet)
	sheetExcept := xlFile.GetSheet(1)
	t.insertToBaseExcept(sheetExcept)

}

func (t *GrlsReader) insertToBase(sheet *xls.WorkSheet) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=OFF&_synchronous=OFF", FileDB))
	if err != nil {
		Logging(err)
		return
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM grls; UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'grls'; VACUUM;")
	if err != nil {
		Logging(err)
		return
	}
	datePub := findFromRegExp(sheet.Row(0).Col(0), `(\d{2}\.\d{2}\.\d{4})`)
	for r := 3; r < int(sheet.MaxRow); r++ {
		col := sheet.Row(r)
		mnn := strings.ReplaceAll(col.Col(0), "\u0000", "")
		mnn = strings.ReplaceAll(mnn, "\u0026", "")
		name := strings.ReplaceAll(col.Col(1), "\u0000", "")
		name = strings.ReplaceAll(name, "\u0000", "")
		form := strings.ReplaceAll(col.Col(2), "\u0000", "")
		form = strings.ReplaceAll(form, "\u0026", "")
		owner := strings.ReplaceAll(col.Col(3), "\u0000", "")
		owner = strings.ReplaceAll(owner, "\u0026", "")
		atx := col.Col(4)
		quantity := col.Col(5)
		maxPrice := strings.ReplaceAll(col.Col(6), ",", ".")
		firstPrice := strings.ReplaceAll(col.Col(7), ",", ".")
		ru := col.Col(8)
		dateReg := col.Col(9)
		code := col.Col(10)
		_, err := db.Exec("INSERT INTO grls (id, mnn, name, form, owner, atx, quantity, max_price, first_price, ru, date_reg, code, date_pub) VALUES (NULL, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", mnn, name, form, owner, atx, quantity, maxPrice, firstPrice, ru, dateReg, code, datePub)
		t.Added++
		if err != nil {
			Logging(err)
		}
	}
}

func (t *GrlsReader) insertToBaseExcept(sheet *xls.WorkSheet) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=OFF&_synchronous=OFF", FileDB))
	if err != nil {
		Logging(err)
		return
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM grls_except; UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'grls_except'; VACUUM;")
	if err != nil {
		Logging(err)
		return
	}
	datePub := findFromRegExp(sheet.Row(0).Col(0), `(\d{2}\.\d{2}\.\d{4})`)
	if datePub == "" {
		Logging("datePub is empty")
	}
	for r := 3; r < int(sheet.MaxRow); r++ {
		col := sheet.Row(r)
		mnn := strings.ReplaceAll(col.Col(0), "\u0000", "")
		mnn = strings.ReplaceAll(mnn, "\u0026", "")
		name := strings.ReplaceAll(col.Col(1), "\u0000", "")
		name = strings.ReplaceAll(name, "\u0000", "")
		form := strings.ReplaceAll(col.Col(2), "\u0000", "")
		form = strings.ReplaceAll(form, "\u0026", "")
		owner := strings.ReplaceAll(col.Col(3), "\u0000", "")
		owner = strings.ReplaceAll(owner, "\u0026", "")
		atx := col.Col(4)
		quantity := col.Col(5)
		maxPrice := strings.ReplaceAll(col.Col(6), ",", ".")
		firstPrice := strings.ReplaceAll(col.Col(7), ",", ".")
		ru := col.Col(8)
		dateReg := col.Col(9)
		code := col.Col(10)
		exceptCause := col.Col(11)
		exceptDate := findFromRegExp(col.Col(13), `(\d{2}\.\d{2}\.\d{4})`)
		if exceptDate == "" {
			Logging(fmt.Sprintf("exceptDate is empty, row %d, mnn - %s", r, mnn))
		}
		_, err := db.Exec("INSERT INTO grls_except (id, mnn, name, form, owner, atx, quantity, max_price, first_price, ru, date_reg, code, except_cause, except_date, date_pub) VALUES (NULL, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", mnn, name, form, owner, atx, quantity, maxPrice, firstPrice, ru, dateReg, code, exceptCause, exceptDate, datePub)
		t.Added++
		if err != nil {
			Logging(err)
		}
	}
}
