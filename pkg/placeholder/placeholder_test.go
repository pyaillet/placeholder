package placeholder

import (
	"fmt"
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

	resultErr := fmt.Errorf("Some values were not found: %v", []string{"THIRD"})

	assert.Nil(t, actual)
	assert.Equal(t, actualErr, resultErr)
}
