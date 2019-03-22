package placeholder

import (
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

	expected := []string{"FIRST", "SECOND"}

	assert.Equal(t, expected, actual)
}
