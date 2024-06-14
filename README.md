# WikiDump Search Engine
a simple full text search engine for wiki abstract dump data

Data Source : https://dumps.wikimedia.org/enwiki/latest/

## Installation

### Prerequisites
Make sure you have the following installed:
- [Go](https://golang.org/doc/install) (version 1.16 or later)
- Git

### Clone the Repository
Clone the repository to your local machine using the following command:

```bash
git clone https://github.com/Haji-sudo/wikidump_searchengine.git
```

Navigate to the project directory:

```bash
cd wikidump_searchengine
```

### Install Dependencies
Use `go mod tidy` to install the required dependencies:

```bash
go mod tidy
```

This command will clean up your `go.mod` file and download the necessary dependencies listed in the `go.mod` file.

### Build the Application
Build the application using the following command:

```bash
go build
```

### Run the Application
Run the application with the following command:

```bash
./appName -file <enwiki-latest-abstract.xml.gz>
```

## Libraries Used
The following libraries are used in this project:

- [`snowball`](github.com/kljensen/snowball/english): Used for stemming English words.
- [`regexp`](https://pkg.go.dev/regexp): Used for find words with a pattern __wildcard__.


