# Setting up the Web app and CLI tool

This guide will take you through the steps on run the web app and cli tool

## Why

This application allows users to retrieve a random photo and quote, making calls to the 3rd party apis https://picsum.photos/ and https://forismatic.com/en/api/ respectively.

The Web app serves an interactive HTML page showing the random photo and quote

The CLI tool provides a back-end interface to also retreive the photo and quote. With the CLI tool, the retrieved photo and quote are stored in a file specified with command line arguments.

Both the Web app and CLI tool [provides options](#supported-options) to manipulate the desired photo and quote to be displayed.

**Note: if a request to the 3rd party provider api fails, this application serves the default [photo](web/img/img-1-gs.jpg) and/or [quote](web/data/quote-1.txt) instead**

## Supported options

- Grayscale - Specifies if the image returned is grayscale or not
- Category - Specifies the cateogry for the quote which influences the choice of the quotation (A numerical value with a max length of 6 characters)

## Pre-requisites
- Go ([follow the guide](https://go.dev/doc/install) to download install the latest version if not already installed)

## 3rd party libraries

- [go-chi](https://github.com/go-chi/chi/) - A lightweight, idiomatic and composable router for building Go HTTP services
- [afero](github.com/spf13/afero) - A filesystem framework for interacting with any filesystem

## Setting up the Web App

### Step 1: Build

Note: The command builds the binaries for the web app in the build/ sub directory

```bash
go build -o build/ cmd/server/*.go
```

### Step 2: Run the Web App

```bash
./build/tucows-exercise-server
```

Go to the following link on your prefered browser: http://localhost:8080/

## CLI Tool

### Step 1: Build

Note: The command builds the binaries for the CLI tool in the build/ sub directory

```bash
go build -o build/ cmd/cli/*.go
```

### Step 2: Run the CLI Tool

```bash
./build/tucows-exercise-cli
```

### Supported command line arguments

Use the --help argument to show all supported arguments and information on how to use the tool

```
./build/tucows-exercise-cli --help

Usage of tucows-exercise-cli:
  --help            Displays information on how to use this tool
  --grayscale       Generates grayscale image
  --category        Specifies the quote category key (default "100000"). This is a numerical value with a maximum length of 6 characters. Note: if an invalid value is specified, the default is used
  --output-image    Specifies the name of the output image file (default "tmp/image.jpg"). Generated file is a jpeg image
  --output-quote    Specifies the name of the output file containing the quote (default "tmp/quote.txt")
```
