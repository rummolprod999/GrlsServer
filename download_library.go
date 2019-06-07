package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func DownloadPage(url string) string {
	count := 0
	var st string
	for {
		//fmt.Println("Start download file")
		if count > 5 {
			Logging(fmt.Sprintf("Не скачали файл за %d попыток %s", count, url))
			return st
		}
		st = GetPage(url)
		if st == "" {
			count++
			Logging("Получили пустую страницу", url)
			time.Sleep(time.Second * 5)
			continue
		}
		return st

	}
	return st
}

func GetPage(url string) string {
	var st string
	resp, err := http.Get(url)
	if err != nil {
		Logging("Ошибка response", url, err)
		return st
	}
	defer resp.Body.Close()
	if err != nil {
		Logging("Ошибка скачивания", url, err)
		return st
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logging("Ошибка чтения", url, err)
		return st
	}

	return string(body)
}

func DownloadF(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func DownloadFile(filepath string, url string) error {
	count := 0
	for {
		if count > 5 {
			return errors.New(fmt.Sprintf("Не скачали файл за %d попыток %s", count, url))
		}
		err := DownloadF(filepath, url)
		if err != nil {
			count++
			Logging(err)
			time.Sleep(time.Second * 5)
			CreateTempDir()
			continue
		}
		return nil

	}
}
