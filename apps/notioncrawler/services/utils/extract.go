package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ExportFileName  = "notion_dump"
	ExportDirectory = "/Users/khadim/dev/notioncrawler/data"
)

func ExtractExportZip(file string) (string, error) {
	now := time.Now()
	outputDirName := fmt.Sprintf("%s-%s",
		ExportFileName,
		now.Format("2006-01-02_15-04-05"),
	)
	outputDir := filepath.Join(ExportDirectory, outputDirName)

	log.Printf("Extracting download: %s", file)

	if err := Unzip(file, outputDir); err != nil {
		return "", err
	}

	if err := os.Remove(file); err != nil {
		log.Printf("Could not delete zip file: %s", file)
	}

	log.Printf("Download extracted: %s", outputDir)
	return outputDir, nil
}

func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		err := UnzipFile(file, dest)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnzipFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	filePath := filepath.Join(dest, f.Name)
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, f.Mode()); err != nil {
			return err
		}
	} else {
		var fileDir string
		if lastIndex := strings.LastIndex(filePath, string(os.PathSeparator)); lastIndex > -1 {
			fileDir = filePath[:lastIndex]
		}

		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
				return err
			}
		}

		f, err := os.OpenFile(
			filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}
