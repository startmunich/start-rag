package utils

import (
	"fmt"
	"log"
	"notioncrawl/services/notion"
	"time"
)

func ExportZip(client *notion.Client, options notion.ExportOptions) (string, error) {
	taskId, err := client.TriggerExportTask(options)
	if err != nil {
		return "", err
	}
	log.Println("taskId extracted:", taskId)

	downloadLink, err := client.GetDownloadLink(taskId)
	if err != nil {
		log.Printf("downloadLink could not be extracted: %v", err)
		return "", err
	}

	log.Println("Download link extracted:", downloadLink)

	log.Println("Downloading file...")
	now := time.Now()
	fileName := fmt.Sprintf("%s-%s-%s.zip",
		ExportFileName,
		options.ExportType,
		now.Format("2006-01-02_15-04-05"),
	)

	log.Println("fileName:", fileName)

	dFile, dPath, err := client.DownloadToFile(downloadLink, fileName)
	log.Printf("--- %s ---", dPath)
	if err != nil {
		log.Printf("Could not download file: %v", err)
		return "", err
	}

	log.Printf("Download finished: %s", dFile.Name())
	return dPath, nil
}

func ExportExtracted(client *notion.Client, options notion.ExportOptions) (string, error) {
	file, err := ExportZip(client, options)
	if err != nil {
		return "", err
	}

	outputDir, err := ExtractExportZip(file)
	if err != nil {
		return "", err
	}

	return outputDir, nil
}
