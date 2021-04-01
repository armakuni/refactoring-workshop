package refactored

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapitalise(t *testing.T) {
	assert.Equal(t, "FOO", Capitalise("foo"))
	assert.Equal(t, "FOO", Capitalise("Foo"))
	assert.Equal(t, "FOO", Capitalise("FOO"))
}

func TestDecapitalisation(t *testing.T) {
	assert.Equal(t, "foo", Decapitalise("foo"))
	assert.Equal(t, "foo", Decapitalise("Foo"))
	assert.Equal(t, "foo", Decapitalise("FOO"))
}

func TestTransformSlice(t *testing.T) {
	input := []string{"Foo", "BAR", "wibble"}
	expected := []string{"FOO", "BAR", "WIBBLE"}
	assert.Equal(t, expected, TransformSlice(input, Capitalise))
}
func TestChooseFormat(t *testing.T) {
	tests := []struct {
		filename    string
		exactFormat string
		expected    string
		errors      bool
	}{
		{"/bar/foo.json", "", JSON, false},
		{"/bar/foo.json", YAML, YAML, false},
		{"/bar/foo.yaml", "", YAML, false},
		{"/bar/foo.yaml", JSON, JSON, false},
		{"/bar/foo.yml", "", YAML, false},
		{"/bar/foo.yml", JSON, JSON, false},
		{"/bar/foo.yml", "fail", "", true},
		{"", "", "", true},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s with %s", test.filename, test.exactFormat), func(t *testing.T) {
			format, err := ChooseFormat(test.filename, test.exactFormat)
			if !test.errors {
				require.NoError(t, err)
				assert.Equal(t, test.expected, format)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func withTempFile(t *testing.T, format string, content string, fn func(*os.File)) {
	file, err := ioutil.TempFile("", format)
	require.NoError(t, err)
	defer os.Remove(file.Name())
	if content != "" {
		if _, err := file.WriteString(content); err != nil {
			t.Fatal(err)
		}
		file.Sync()
	}
	fn(file)
}

func readFile(t *testing.T, fp *os.File) string {
	data, err := ioutil.ReadAll(fp)
	require.NoError(t, err)
	return string(data)
}

func TestWithReaderFile(t *testing.T) {
	withTempFile(t, "*.txt", "test data", func(file *os.File) {
		err := WithReader(file.Name(), nil, func(r io.Reader) error {
			data, err2 := ioutil.ReadAll(r)
			require.NoError(t, err2)
			assert.Equal(t, "test data", string(data))
			return nil
		})
		assert.NoError(t, err)
	})
}

func TestWithReaderStream(t *testing.T) {
	buf := strings.NewReader("test data")
	err := WithReader("", buf, func(r io.Reader) error {
		data, err2 := ioutil.ReadAll(r)
		require.NoError(t, err2)
		assert.Equal(t, "test data", string(data))
		return nil
	})
	assert.NoError(t, err)
}

func TestWithWriterFile(t *testing.T) {
	withTempFile(t, "*.txt", "", func(file *os.File) {
		err := WithWriter(file.Name(), nil, func(w io.Writer) error {
			_, err2 := w.Write([]byte("test data"))
			require.NoError(t, err2)
			return nil
		})
		require.NoError(t, err)
		data, err := ioutil.ReadAll(file)
		require.NoError(t, err)
		assert.Equal(t, "test data", string(data))
	})
}

func TestWithWriterStream(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	err := WithWriter("", buf, func(w io.Writer) error {
		_, err2 := w.Write([]byte("test data"))
		require.NoError(t, err2)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "test data", buf.String())
}

func TestCreateDecoder(t *testing.T) {
	r := strings.NewReader("foo")
	assert.NotNil(t, CreateDecoder(r, JSON))
	assert.NotNil(t, CreateDecoder(r, YAML))
	assert.Nil(t, CreateDecoder(r, ""))
}

func TestCreateEncoder(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	assert.NotNil(t, CreateEncoder(buf, JSON))
	assert.NotNil(t, CreateEncoder(buf, YAML))
	assert.Nil(t, CreateEncoder(buf, ""))
}
