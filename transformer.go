package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	var filename string
	var stream io.Reader
	var inputFormat string
	var outputFormat string
	var transformation string
	var outputFile string
	var outputStream io.Writer

	flag.StringVar(&filename, "i", "", "File to read from, or '-' for STDIN")
	flag.StringVar(&inputFormat, "f", "", "Input format - JSON or YAML, if blank will be derived from filename")
	flag.StringVar(&outputFormat, "F", "", "Output format - JSON or YAML, if blank will be derived from filename")
	flag.StringVar(&transformation, "t", "CAPITALISE", "Transformation to apply")
	flag.StringVar(&outputFile, "o", "", "File to write to, or '-' for STDOUT")

	flag.Parse()

	if filename == "-" {
		stream = os.Stdin
		filename = ""
	}

	if outputFile == "-" {
		outputStream = os.Stdout
		outputFile = ""
	}

	if err := Transform(filename, stream, inputFormat, outputFormat, transformation, outputFile, outputStream); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func Transform(filename string, stream io.Reader, inputFormat, outputFormat, transformation, outputFile string, outputStream io.Writer) (err error) {
	stuff := make([]string, 0)
	var raw_stuff []byte

	if filename != "" {
		fp, err := os.Open(filename)
		if err != nil {
			return err
		}
		if inputFormat == "" {
			if strings.HasSuffix(filename, ".json") {
				inputFormat = "JSON"
			} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
				inputFormat = "YAML"
			}
			if inputFormat == "" {
				fp.Close()
				return errors.New("No input format was specified")
			}
		}
		if inputFormat == "YAML" {
			raw_stuff, err = ioutil.ReadAll(fp)
			if err == nil {
				err = yaml.Unmarshal(raw_stuff, &stuff)
			}
		} else {
			decoder := json.NewDecoder(fp)
			if err := decoder.Decode(&stuff); err != nil {
				return err
			}
		}
		fp.Close()
		if err != nil {
			return err
		}
	} else {
		raw_stuff, err = ioutil.ReadAll(stream)
		if inputFormat == "" {
			return errors.New("No input format was specified")
		}
		if inputFormat == "YAML" {
			if err := yaml.Unmarshal(raw_stuff, &stuff); err != nil {
				return err
			}
		} else if inputFormat == "JSON" {
			if err := json.Unmarshal(raw_stuff, &stuff); err != nil {
				return err
			}
		}
	}

	if outputFormat == "" {
		if strings.HasSuffix(outputFile, ".json") {
			outputFormat = "JSON"
		} else if strings.HasSuffix(outputFile, ".yaml") || strings.HasSuffix(outputFile, ".yml") {
			outputFormat = "YAML"
		} else {
			outputFormat = inputFormat
		}
	}
	more_stuff := make([]string, 0)

	if transformation == "CAPITALISE" {
		for _, element := range stuff {
			more_stuff = append(more_stuff, strings.ToUpper(element))
		}
	} else if transformation == "DECAPITALISE" {
		for _, element := range stuff {
			more_stuff = append(more_stuff, strings.ToLower(element))
		}
	}

	if outputFormat == "JSON" {
		if outputFile != "" {
			data, err := json.Marshal(&more_stuff)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(outputFile, data, 0666); err != nil {
				return err
			}
		} else if outputStream != nil {
			encoder := json.NewEncoder(outputStream)
			encoder.Encode(&more_stuff)
		} else {
			return errors.New("No output stream is specified")
		}
	}

	if outputFormat == "YAML" {
		data, err := yaml.Marshal(&more_stuff)
		if err != nil {
			return err
		}
		if outputStream != nil {
			if _, err := outputStream.Write(data); err != nil {
				return err
			}
		} else {
			if err := ioutil.WriteFile(outputFile, data, 0666); err != nil {
				return err
			}
		}
	}

	return nil
}
