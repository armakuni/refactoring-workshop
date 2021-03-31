package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var contentJSON = `["Foo", "BAR", "wibble"]`
var uppercaseJSON = `["FOO","BAR","WIBBLE"]`
var lowercaseJSON = `["foo","bar","wibble"]`
var contentYAML = `- Foo
- BAR
- wibble
`
var uppercaseYAML = `- FOO
- BAR
- WIBBLE
`
var lowercaseYAML = `- foo
- bar
- wibble
`

func tempFile(t *testing.T, format string, content string) *os.File {
	file, err := ioutil.TempFile("", format)
	require.NoError(t, err)
	if content != "" {
		_, err := file.WriteString(content)
		if err != nil {
			os.Remove(file.Name())
			t.Fatal(err)
		}
		file.Sync()
	}
	return file
}
func TestTransformerJSONFileToJSONFileCAPITALISE(t *testing.T) {
	inputFile := tempFile(t, "*.json", contentJSON)
	defer os.Remove(inputFile.Name())
	outputFile := tempFile(t, "*.json", "")
	defer os.Remove(outputFile.Name())

	err := Transform(inputFile.Name(), nil, "", "", "CAPITALISE", outputFile.Name(), nil)
	require.NoError(t, err)
	content, err := ioutil.ReadAll(outputFile)
	require.NoError(t, err)
	assert.Equal(t, uppercaseJSON, string(content))
}

func TestTransformerJSONStreamToJSONStreamDECAPITALISE(t *testing.T) {
	inputStream := strings.NewReader(contentJSON)
	buf := bytes.NewBufferString("")

	err := Transform("", inputStream, "JSON", "JSON", "DECAPITALISE", "", buf)
	require.NoError(t, err)
	assert.Equal(t, lowercaseJSON, strings.TrimSpace(buf.String()))
}

func TestTransformerYAMLFileToYAMLFileCAPITALISE(t *testing.T) {
	inputFile := tempFile(t, "*.yaml", contentYAML)
	defer os.Remove(inputFile.Name())
	outputFile := tempFile(t, "*.yaml", "")
	defer os.Remove(outputFile.Name())

	err := Transform(inputFile.Name(), nil, "", "", "CAPITALISE", outputFile.Name(), nil)
	require.NoError(t, err)
	content, err := ioutil.ReadAll(outputFile)
	require.NoError(t, err)
	assert.Equal(t, uppercaseYAML, string(content))
}

func TestTransformerYAMLStreamToJSONStreamDECAPITALISE(t *testing.T) {
	inputStream := strings.NewReader(contentJSON)
	buf := bytes.NewBufferString("")

	err := Transform("", inputStream, "YAML", "YAML", "DECAPITALISE", "", buf)
	require.NoError(t, err)
	assert.Equal(t, lowercaseYAML, buf.String())
}

// TODO:
//
// JSON stream to JSON file, CAPITALISE & DECAPITALISE
// JSON file to JSON stream, CAPITALISE & DECAPITALISE
// JSON file to YAML file, CAPITALISE & DECAPITALISE
// JSON file to YAML stream, CAPITALISE & DECAPITALISE
// YAML file to YAML stream, CAPITALISE & DECAPITALISE
// YAML stream to YAML file, CAPITALISE & DECAPITALISE
// YAML stream to JSON file/stream, CAPITALISE & DECAPITALISE
// edge cases such as inferring file format, specifying file format, specifying the wrong file format, general errors
// cartesion product of the above, to be absolutely sure
