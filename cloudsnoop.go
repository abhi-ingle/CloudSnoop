package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"code.sajari.com/docconv/v2"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Define command-line flags for URLs and snoop options
	urlFlag := flag.String("url", "", "Comma-separated URLs to fetch")
	fileFlag := flag.String("file", "", "Path to a file containing URLs")
	snoopFlag := flag.String("snoop", "drive", "Types of links to extract (drive, sharepoint, dropbox, all)")
	flag.Parse()

	// Define the dump folder path
	dumpFolder := "dump"

	// Create the dump folder if it doesn't exist
	if err := os.MkdirAll(dumpFolder, os.ModePerm); err != nil {
		fmt.Printf("Failed to create dump folder: %s\n", err)
		return
	}

	// Ensure the dump folder is removed after program execution
	defer func() {
		os.RemoveAll(dumpFolder)
	}()

	var urls []string

	// Check if a file is provided
	if *fileFlag != "" {
		content, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Printf("Failed to read file: %s\n", err)
			return
		}
		urls = strings.Split(string(content), "\n")
	} else if *urlFlag != "" {
		// Split the URLs into a slice
		urls = strings.Split(*urlFlag, ",")
	} else {
		fmt.Println("Please provide URLs using the --url flag or a file using the --file flag.")
		return
	}

	// Loop through each URL
	for _, url := range urls {
		url = strings.TrimSpace(url) // Trim whitespace
		if url == "" {
			continue // Skip empty lines
		}
		fmt.Printf("Processing URL: %s\n", url)

		// Download the file to disk in the dump folder
		filePath, err := downloadFileToDisk(url, dumpFolder)
		if err != nil {
			fmt.Printf("Failed to download the file: %s\n", err)
			continue
		}

		// Extract text from the downloaded file
		extractedContent, err := extractTextFromFile(url, filePath)
		if err != nil {
			fmt.Printf("Failed to extract text from the file: %s\n", err)
			continue
		}

		// Initialize a slice for extracted links
		var driveLinks, sharepointLinks, dropboxLinks []string

		// Regular expressions for different services
		if shouldExtract("drive", *snoopFlag) || *snoopFlag == "all" {
			driveLinkRegex := regexp.MustCompile(`https?://drive\.google\.com/[^\s"'>]+`)
			driveLinks = driveLinkRegex.FindAllString(extractedContent, -1)
		}
		if shouldExtract("sharepoint", *snoopFlag) || *snoopFlag == "all" {
			sharepointLinkRegex := regexp.MustCompile(`https?://([a-z0-9\-]+\.|)(my\.|team\.|)[a-z0-9\-]+\.sharepoint\.com/[^\s"'>]+`)
			sharepointLinks = sharepointLinkRegex.FindAllString(extractedContent, -1)
		}
		if shouldExtract("dropbox", *snoopFlag) || *snoopFlag == "all" {
			dropboxLinkRegex := regexp.MustCompile(`https?://[^\s"'>]+dropbox\.com/[^\s"'>]+`)
			dropboxLinks = dropboxLinkRegex.FindAllString(extractedContent, -1)
		}

		// Print found links
		printLinks("Google Drive links", driveLinks)
		printLinks("SharePoint links", sharepointLinks)
		printLinks("Dropbox links", dropboxLinks)
	}
}

// Function to check if the current service should be extracted
func shouldExtract(service string, snoop string) bool {
	services := strings.Split(snoop, ",")
	for _, s := range services {
		if strings.TrimSpace(s) == service {
			return true
		}
	}
	return false
}

// Function to print found links
func printLinks(serviceName string, links []string) {
	if len(links) == 0 {
		return
	}
	fmt.Printf("Found %s:\n", serviceName)
	for _, link := range links {
		fmt.Println(link)
	}
}

// Function to download the file to disk in the specified dump folder
func downloadFileToDisk(url, folder string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Get the filename from the URL
	fileName := filepath.Base(url)
	if fileName == "" || fileName == "." || fileName == "/" {
		fileName = "downloaded_file"
	}

	// Create the full path for the downloaded file
	filePath := filepath.Join(folder, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// Function to extract text from the downloaded file
func extractTextFromFile(url, filePath string) (string, error) {
	// Check file extension
	if strings.HasSuffix(url, ".txt") {
		// Read and return text file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	if strings.HasSuffix(url, ".pdf") || strings.HasSuffix(url, ".pptx") {
		// Use docconv to extract text from the file
		res, err := docconv.ConvertPath(filePath)
		if err != nil {
			return "", err
		}
		return res.Body, nil
	}

	// If it's not a txt, pdf, or pptx, treat it as HTML and scrape
	return extractTextFromHTMLFile(filePath), nil
}

// Function to extract text from HTML file content using goquery
func extractTextFromHTMLFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return "" // Return empty string on error
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return "" // Return empty string on error
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(content)))
	if err != nil {
		return "" // Return empty string on error
	}
	// Get the text from the document
	return doc.Find("body").Text()
}
