package refactored

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	JSON         = "JSON"
	YAML         = "YAML"
	CAPITALISE   = "CAPITALISE"
	DECAPITALISE = "DECAPITALISE"
)

func Transform(inputFilename, inputFormat, outputFilename, outputFormat, transformation string) error {
	var r io.Reader
	var w io.Writer
	if inputFilename == "-" {
		r = os.Stdin
		inputFilename = ""
	}
	if outputFilename == "-" {
		w = os.Stdout
		outputFilename = ""
	}
	transformer, err := ChooseTransformer(transformation)
	if err != nil {
		return err
	}
	inputFormat, err = ChooseFormat(inputFilename, inputFormat)
	if err != nil {
		return err
	}
	outputFormat, err = ChooseFormat(outputFilename, outputFormat)
	if err != nil {
		return err
	}
	return WithReader(inputFilename, r, func(r io.Reader) error {
		return WithWriter(outputFilename, w, func(w io.Writer) error {
			decoder := CreateDecoder(r, inputFormat)
			encoder := CreateEncoder(w, outputFormat)
			return DecodeTransformAndEncode(decoder, encoder, transformer)
		})
	})
}

func DecodeTransformAndEncode(decoder Decoder, encoder Encoder, transformer TransformerFunc) error {
	var elements []string
	if err := decoder.Decode(&elements); err != nil {
		return err
	}
	transformedElements := TransformSlice(elements, transformer)
	return encoder.Encode(&transformedElements)
}

type TransformerFunc func(string) string

func Capitalise(input string) string {
	return strings.ToUpper(input)
}

func Decapitalise(input string) string {
	return strings.ToLower(input)
}

func TransformSlice(input []string, transformer TransformerFunc) []string {
	output := make([]string, len(input))
	for idx, element := range input {
		output[idx] = transformer(element)
	}
	return output
}

func ChooseTransformer(transformer string) (TransformerFunc, error) {
	switch transformer {
	case CAPITALISE:
		return Capitalise, nil
	case DECAPITALISE:
		return Decapitalise, nil
	}
	return nil, fmt.Errorf("Cannot determine transformation function from '%s'", transformer)
}

func ChooseFormat(filename string, exactFormat string) (format string, err error) {
	format = exactFormat
	if format == "" {
		switch filepath.Ext(filename) {
		case ".json":
			format = JSON
		case ".yaml":
			format = YAML
		case ".yml":
			format = YAML
		}
	}
	switch format {
	case JSON:
	case YAML:
	case "":
		err = errors.New("Cannot determine format")
	default:
		err = fmt.Errorf("Incorrect format specified: '%s'", format)
	}
	return
}

type Encoder interface {
	Encode(v interface{}) error
}

type Decoder interface {
	Decode(v interface{}) error
}

func CreateDecoder(r io.Reader, format string) Decoder {
	switch format {
	case JSON:
		return json.NewDecoder(r)
	case YAML:
		return yaml.NewDecoder(r)
	}
	return nil
}

func CreateEncoder(w io.Writer, format string) Encoder {
	switch format {
	case JSON:
		return json.NewEncoder(w)
	case YAML:
		return yaml.NewEncoder(w)
	}
	return nil
}

func WithReader(filename string, stream io.Reader, fn func(io.Reader) error) error {
	if stream != nil {
		return fn(stream)
	}
	if fp, err := os.Open(filename); err != nil {
		return err
	} else {
		defer fp.Close()
		return fn(fp)
	}
}

func WithWriter(filename string, stream io.Writer, fn func(io.Writer) error) error {
	if stream != nil {
		return fn(stream)
	}
	if fp, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		return err
	} else {
		defer fp.Close()
		return fn(fp)
	}
}
