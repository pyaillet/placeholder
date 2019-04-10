package placeholder

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPlaceHolders_ok(t *testing.T) {
	data := "%#FIRST#% azpoerkjapzoje \n %#SECOND#% %#THIRD#%"
	sep := Separator{
		start: "%#",
		end:   "#%",
	}

	actual := ListPlaceHolders([]byte(data), sep)

	expected := []string{"FIRST", "SECOND", "THIRD"}

	assert.Equal(t, expected, actual)

}

func TestListPlaceHolders_ok_in_order(t *testing.T) {
	data := "%#THIRD#% %#FIRST#% azpoerkjapzoje \n %#SECOND#% %#THIRD#%"
	sep := Separator{
		start: "%#",
		end:   "#%",
	}

	actual := ListPlaceHolders([]byte(data), sep)

	expected := []string{"FIRST", "SECOND", "THIRD"}

	assert.Equal(t, expected, actual)

}

func TestListPlaceHolders_alternate_ok_in_order(t *testing.T) {
	data := "{{THIRD}} {{FIRST}} azpoerkjapzoje \n {{SECOND}} {{THIRD}}"
	sep := SeparatorFrom("{{", "}}")

	actual := ListPlaceHolders([]byte(data), sep)

	expected := []string{"FIRST", "SECOND", "THIRD"}

	assert.Equal(t, expected, actual)

}

func TestListPlaceHolders_uniq(t *testing.T) {
	data := "%#FIRST#%\n%#FIRST#%"
	sep := Separator{
		start: "%#",
		end:   "#%",
	}

	actual := ListPlaceHolders([]byte(data), sep)

	expected := []string{"FIRST"}

	assert.Equal(t, expected, actual)

}

func TestUniq_uniq(t *testing.T) {
	data := []string{"FIRST", "FIRST"}

	actual := uniq(data)

	expected := []string{"FIRST"}

	assert.Equal(t, expected, actual)
}

func TestUniq_one(t *testing.T) {
	data := []string{"FIRST"}

	actual := uniq(data)

	expected := []string{"FIRST"}

	assert.Equal(t, expected, actual)
}

func TestUniq_two(t *testing.T) {
	data := []string{"FIRST", "SECOND"}

	actual := uniq(data)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, "FIRST")
	assert.Contains(t, actual, "SECOND")
}

func TestListPlaceHoldersInFiles_ok(t *testing.T) {
	data := []string{"./testdata/example.html", "./testdata/example.js"}

	actual := ListPlaceHoldersInFiles(data, DefaultSeparator())

	expected := []string{"INDEX", "MESSAGE_WITH_COMPOSED_KEY", "TITLE"}

	assert.Equal(t, expected, actual)
}

func TestReplacingPlaceHolders_ok(t *testing.T) {
	data := "%#FIRST#% azpoerkjapzoje %#THIRD#% \n %#SECOND#%"

	values := map[string]string{
		"FIRST":  "first",
		"SECOND": "second",
		"THIRD":  "third",
	}

	actual := ReplacingPlaceHolders([]byte(data), values, DefaultSeparator())

	expected := "first azpoerkjapzoje third \n second"

	assert.Equal(t, []byte(expected), actual)
}

func TestReplacingPlaceHoldersFromEnv_ok(t *testing.T) {
	data := "%#FIRST#% azpoerkjapzoje %#THIRD#% \n %#SECOND#%"

	os.Setenv("FIRST", "first")
	os.Setenv("SECOND", "second")
	os.Setenv("THIRD", "third")

	actual, err := ReplacingPlaceHoldersFromEnv([]byte(data), DefaultSeparator())

	expected := "first azpoerkjapzoje third \n second"

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), actual)
}

func TestReplacingPlaceHoldersFromEnv_error_some_undefined(t *testing.T) {
	data := "%#FIRST#% azpoerkjapzoje %#THIRD#% \n %#SECOND#%"

	os.Setenv("FIRST", "first")
	os.Setenv("SECOND", "second")
	os.Unsetenv("THIRD")

	actual, actualErr := ReplacingPlaceHoldersFromEnv([]byte(data), DefaultSeparator())

	resultErr := fmt.Errorf("Some values were not found: %+q", []string{"THIRD"})

	assert.Nil(t, actual)
	assert.Equal(t, actualErr, resultErr)
}

func TestReplacingPlaceHoldersInFilesFromEnv_ok(t *testing.T) {
	copyFile("./testdata/example.html", "./testdata/example1.html")
	copyFile("./testdata/example.js", "./testdata/example1.js")

	os.Setenv("TITLE", "My title")
	os.Setenv("MESSAGE_WITH_COMPOSED_KEY", "This is my message")
	os.Setenv("INDEX", "0")

	files := []string{"./testdata/example1.html", "./testdata/example1.js"}
	err := ReplacingPlaceHoldersInFilesFromEnv(files, DefaultSeparator())

	assert.Nil(t, err)
	assertSameFileContent(t, "./testdata/example1.html", "./testdata/example.html_result")
	assertSameFileContent(t, "./testdata/example1.js", "./testdata/example.js_result")

	os.Remove("./testdata/example1.html")
	os.Remove("./testdata/example1.js")
}

func TestReplacingPlaceHoldersInFilesFromEnv_ko_some_values_not_found(t *testing.T) {
	copyFile("./testdata/example.html", "./testdata/example1.html")
	copyFile("./testdata/example.js", "./testdata/example1.js")

	os.Setenv("TITLE", "My title")
	os.Unsetenv("MESSAGE_WITH_COMPOSED_KEY")
	os.Setenv("INDEX", "0")

	files := []string{"./testdata/example1.html", "./testdata/example1.js"}
	actualErr := ReplacingPlaceHoldersInFilesFromEnv(files, DefaultSeparator())

	expectedErr := fmt.Errorf("Some values were not found: %+q", []string{"MESSAGE_WITH_COMPOSED_KEY"})

	assert.Equal(t, actualErr, expectedErr)

	os.Remove("./testdata/example1.html")
	os.Remove("./testdata/example1.js")
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func assertSameFileContent(t *testing.T, expected, actual string) {
	actualContent, err := ioutil.ReadFile(actual)
	assert.Nil(t, err)

	expectedContent, err := ioutil.ReadFile(expected)
	assert.Nil(t, err)

	assert.Equal(t, expectedContent, actualContent)
}
