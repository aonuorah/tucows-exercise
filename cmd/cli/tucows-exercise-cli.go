package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aonuorah/tucows-exercise/internal/photo"
	"github.com/aonuorah/tucows-exercise/internal/quote"
	"github.com/aonuorah/tucows-exercise/internal/service"
)

func main() {
	args := parseArgs()
	ctx := args.WithContext(context.TODO())
	photo, quote := service.GetNewItem(ctx)

	if err := service.WriteToFile(photo.Bytes, args.ImageFile); err != nil {
		panic(err)
	}

	if err := service.WriteToFile([]byte(quote.Text), args.QuoteFile); err != nil {
		panic(err)
	}

	fmt.Println("Photo: " + args.ImageFile)
	fmt.Println("Quote: " + args.QuoteFile)
}

const usage = `
Usage of tucows-exercise-cli:
  --help 	Displays information on how to use this tool
  --grayscale 	Generates grayscale image
  --category 	Specifies the quote category key (default "100000"). This is a numerical value with a maximum length of 6 characters
	Note: if an invalid value is specified, the default is used
  --output-image 	Specifies the name of the output image file (default "tmp/image.jpg"). Generated file is a jpeg image
  --output-quote 	Specifies the name of the output file containing the quote (default "tmp/quote.txt")`

func parseArgs() *args {
	var a args
	var help bool

	f := flag.NewFlagSet("tucows-exercise", flag.ExitOnError)

	f.BoolVar(&help, "help", false, "Displays information on how to use this tool")
	f.BoolVar(&a.Grayscale, "grayscale", false, "Generates grayscale image")
	f.IntVar(&a.Category, "category", 100000, "Specifies the quote category key")
	f.StringVar(&a.ImageFile, "output-image", "tmp/image.jpg", "Specifies the name of the output image file")
	f.StringVar(&a.QuoteFile, "output-quote", "tmp/quote.txt", "Specifies the name of the output quote file")

	f.Usage = func() { fmt.Println(usage) }
	f.Parse(os.Args[1:])

	if help {
		f.Usage()
		os.Exit(0)
	}

	if err := a.Valid(); err != nil {
		fmt.Println(err.Error())
		f.Usage()
		os.Exit(2)
	}

	return &a
}

// struct to hold the values of all supported cli arguments passed to the application
type args struct {
	Grayscale bool
	Category  int
	ImageFile string
	QuoteFile string
}

// Validates the argument values returning an error if any validation fails
func (a args) Valid() error {
	if a.ImageFile == "" {
		return errors.New("'output-image' is required")
	}

	if a.QuoteFile == "" {
		return errors.New("'output-quote' is required")
	}

	return nil
}

// This method takes a context and returns a new context with the required argument value set
func (a args) WithContext(ctx context.Context) context.Context {
	ctx = photo.CtxWithGrayscale(ctx, a.Grayscale)
	ctx = quote.CtxWithCategory(ctx, a.Category)
	return ctx
}
