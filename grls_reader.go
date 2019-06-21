package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/extrame/xls"
	"os"
	"path/filepath"
	"strings"
)

type GrlsReader struct {
	Url         string
	Added       int
	AddedExcept int
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
	db, err := DbConnection()
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
	datePub := FindFromRegExp(sheet.Row(0).Col(0), `(\d{2}\.\d{2}\.\d{4})`)
	for r := 3; r <= int(sheet.MaxRow); r++ {
		col := sheet.Row(r)
		mnn := ReplaceBadSymbols(col.Col(0))
		name := ReplaceBadSymbols(col.Col(1))
		form := ReplaceBadSymbols(col.Col(2))
		owner := ReplaceBadSymbols(col.Col(3))
		atx := ReplaceBadSymbols(col.Col(4))
		quantity := ReplaceBadSymbols(col.Col(5))
		maxPrice := strings.ReplaceAll(ReplaceBadSymbols(col.Col(6)), ",", ".")
		firstPrice := strings.ReplaceAll(ReplaceBadSymbols(col.Col(7)), ",", ".")
		ru := ReplaceBadSymbols(col.Col(8))
		dateReg := ReplaceBadSymbols(col.Col(9))
		code := ReplaceBadSymbols(col.Col(10))
		if mnn == "" && name == "" && form == "" && owner == "" && atx == "" && quantity == "" && maxPrice == "" && firstPrice == "" && ru == "" && code == "" {
			return
		}
		_, err := db.Exec("INSERT INTO grls (id, mnn, name, form, owner, atx, quantity, max_price, first_price, ru, date_reg, code, date_pub) VALUES (NULL, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", mnn, name, form, owner, atx, quantity, maxPrice, firstPrice, ru, dateReg, code, datePub)
		t.Added++
		if err != nil {
			Logging(err)
		}
	}
}

func (t *GrlsReader) insertToBaseExcept(sheet *xls.WorkSheet) {
	db, err := DbConnection()
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
	datePub := FindFromRegExp(sheet.Row(0).Col(0), `(\d{2}\.\d{2}\.\d{4})`)
	if datePub == "" {
		Logging("datePub is empty")
	}
	for r := 3; r <= int(sheet.MaxRow); r++ {
		col := sheet.Row(r)
		mnn := ReplaceBadSymbols(col.Col(0))
		name := ReplaceBadSymbols(col.Col(1))
		form := ReplaceBadSymbols(col.Col(2))
		owner := ReplaceBadSymbols(col.Col(3))
		atx := ReplaceBadSymbols(col.Col(4))
		quantity := ReplaceBadSymbols(col.Col(5))
		maxPrice := strings.ReplaceAll(ReplaceBadSymbols(col.Col(6)), ",", ".")
		firstPrice := strings.ReplaceAll(ReplaceBadSymbols(col.Col(7)), ",", ".")
		ru := ReplaceBadSymbols(col.Col(8))
		dateReg := ReplaceBadSymbols(col.Col(9))
		code := ReplaceBadSymbols(col.Col(10))
		exceptCause := ReplaceBadSymbols(col.Col(11))
		exceptDate := FindFromRegExp(ReplaceBadSymbols(col.Col(13)), `(\d{2}\.\d{2}\.\d{4})`)
		if exceptDate == "" {
			Logging(fmt.Sprintf("exceptDate is empty, row %d, mnn - %s", r, mnn))
		}
		if mnn == "" && name == "" && form == "" && owner == "" && atx == "" && quantity == "" && maxPrice == "" && firstPrice == "" && ru == "" && code == "" && exceptCause == "" && exceptDate == "" {
			return
		}
		_, err := db.Exec("INSERT INTO grls_except (id, mnn, name, form, owner, atx, quantity, max_price, first_price, ru, date_reg, code, except_cause, except_date, date_pub) VALUES (NULL, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", mnn, name, form, owner, atx, quantity, maxPrice, firstPrice, ru, dateReg, code, exceptCause, exceptDate, datePub)
		t.AddedExcept++
		if err != nil {
			Logging(err)
		}
	}
}
