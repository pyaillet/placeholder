package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var binaryName = "placeholder"

var testPath = "integration/testdata/"

func TestCliListDefaultSeparator(t *testing.T) {
	actual, err := launchToolWithFlags(t,
		"ls", testPath+"default_separator.js")
	assert.Nil(t, err)

	assert.Equal(t, "INDEX\nMAX_VALUE3", string(actual))
}

func TestCliListCustomSeparator(t *testing.T) {
	actual, err := launchToolWithFlags(t,
		"--start", "{{", "--end", "}}", "ls", testPath+"custom_separator.js")
	assert.Nil(t, err)

	assert.Equal(t, "INDEX\nMAX_VALUE3", string(actual))
}

func launchToolWithFlags(t *testing.T, args ...string) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(path.Join(dir+"/dist", binaryName), args...)

	return cmd.CombinedOutput()
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

	assert.Equal(t, string(expectedContent), string(actualContent))
}

func TestReplaceWithEnv(t *testing.T) {
	os.Setenv("INDEX", "1")
	os.Setenv("MAX_VALUE3", "20")

	copyFile(testPath+"default_separator.js", testPath+"default_separator_copy.js")

	_, err := launchToolWithFlags(t, "rp", testPath+"default_separator_copy.js")
	assert.Nil(t, err)

	assertSameFileContent(t, testPath+"default_separator.js_result", testPath+"default_separator_copy.js")
	os.Remove(testPath + "default_separator_copy.js")
}

func TestReplaceWithJsonFile(t *testing.T) {
	copyFile(testPath+"default_separator.js", testPath+"default_separator_copy.js")

	_, err := launchToolWithFlags(t, "rp", "-i", testPath+"values.json", testPath+"default_separator_copy.js")
	assert.Nil(t, err)

	assertSameFileContent(t, testPath+"default_separator.js_result", testPath+"default_separator_copy.js")
	os.Remove(testPath + "default_separator_copy.js")
}

func TestReplaceWithYamlFile(t *testing.T) {
	copyFile(testPath+"default_separator.js", testPath+"default_separator_copy.js")

	_, err := launchToolWithFlags(t, "rp", "-i", testPath+"values.yaml", testPath+"default_separator_copy.js")
	assert.Nil(t, err)

	assertSameFileContent(t, testPath+"default_separator.js_result", testPath+"default_separator_copy.js")
	os.Remove(testPath + "default_separator_copy.js")
}

func TestReplaceWithPropertiesFile(t *testing.T) {
	copyFile(testPath+"default_separator.js", testPath+"default_separator_copy.js")

	_, err := launchToolWithFlags(t, "rp", "-i", testPath+"values.properties", testPath+"default_separator_copy.js")
	assert.Nil(t, err)

	assertSameFileContent(t, testPath+"default_separator.js_result", testPath+"default_separator_copy.js")
	os.Remove(testPath + "default_separator_copy.js")
}

func TestReplaceCustomWithPropertiesFile(t *testing.T) {
	copyFile(testPath+"custom_separator_2.js", testPath+"custom_separator_2_copy.js")

	_, err := launchToolWithFlags(t, "-s", "${", "-e", "}", "rp", "-i", testPath+"values.properties", testPath+"custom_separator_2_copy.js")
	assert.Nil(t, err)

	assertSameFileContent(t, testPath+"default_separator.js_result", testPath+"custom_separator_2_copy.js")
	os.Remove(testPath + "custom_separator_2_copy.js")
}

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}
	make := exec.Command("make", "placeholder")
	err = make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
