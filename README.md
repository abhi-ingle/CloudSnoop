<h1 align="center">
  <br>
  <img src="static/snooper.jpeg" width="80%" height="270px" alt="Snooper">
</h1>

<h1 align="center"><u>Snooper</u></h1>

<p align="center">
  <a href="#description">Description</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#future-scope">Future Scope</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

<hr>

## Description

Snooper is a Go-based tool designed to extract and analyze cloud storage links (Google Drive, SharePoint, Dropbox) from various sources, including URL endpoints and publicly exposed files (`.pptx`, `.pdf`, `.txt`). 

## Installation

1. **Clone the repository**:
    ```bash
    git clone https://github.com/abhi-ingle/Snooper.git
    ```

2. **Navigate to the project directory**:
    ```bash
    cd Snooper
    ```

3. **Build the tool**:
    ```bash
    go build -o snooper snooper.go
    ```

## Usage

Snooper can extract cloud storage links from either directly provided URLs or from a file containing multiple URLs.

### Command-line Options

- `--url`: Comma-separated URLs to process (enclose in quotes).
- `--file`: Path to a file containing URLs (one URL per line).
- `--snoop`: Types of links to extract (`drive`, `sharepoint`, `dropbox`, or `all`).

### Examples

1. **Extract Dropbox links from a file**:
    ```bash
    ./snooper --snoop dropbox --file path/to/your/file.txt
    ```

2. **Extract all cloud links from given URLs**:
    ```bash
    ./snooper --snoop all --url "https://example.com/file1.pdf","https://example.com/file2.pptx"
    ```

3. **Extract Google Drive links from a specific URL**:
    ```bash
    ./snooper --snoop drive --url "https://example.com/public-file.html"
    ```

## Future Scope

Planned features include:

- **Support for More File Types**: Adding support for formats like `.docx`, `.xlsx`, and others.
  
- **Crawling Support**: Enabling extraction of links from nested pages on websites for more comprehensive analysis.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any changes or enhancements you'd like to see. Suggestions for new features and improvements are highly encouraged.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
