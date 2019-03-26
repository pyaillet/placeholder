package placeholder

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Separator structure with start and end delimiters
type Separator struct {
	start string
	end   string
}

func DefaultSeparator() Separator {
	return Separator{
		start: "%#",
		end:   "#%",
	}
}

func postProcess(l []string) []string {
	return sortInPlace(uniq(l))
}

// ListPlaceHolders contained in the provided string
func ListPlaceHolders(data []byte, sep Separator) []string {
	return postProcess(listPlaceHolders(data, sep))
}

func listPlaceHolders(data []byte, sep Separator) []string {
	var l []string
	re := regexp.MustCompile(sep.start + "([A-Za-z][A-Za-z0-9_]*?)" + sep.end)

	matches := re.FindAllSubmatch(data, -1)
	for _, match := range matches {
		l = append(l, string(match[1]))
	}
	return l
}

// ListPlaceHoldersInFiles lists placeholders contained in files
func ListPlaceHoldersInFiles(files []string, sep Separator) []string {
	var all []string
	for _, f := range files {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Errorf("Couldn't read file %s\n", f)
		}
		ph1 := listPlaceHolders(file, sep)
		all = append(all, ph1...)
	}
	return postProcess(all)
}

func sortInPlace(l []string) []string {
	r := make([]string, len(l))
	copy(r, l)
	sort.Strings(r)
	return r
}

func uniq(list []string) []string {
	set := map[string]interface{}{}
	for _, val := range list {
		set[val] = nil
	}
	l := make([]string, 0, len(set))
	for k := range set {
		l = append(l, k)
	}
	return l
}

// ReplacingPlaceHolders replaces place holders in the string according to the
// values map content
func ReplacingPlaceHolders(data []byte, values map[string]string, sep Separator) []byte {
	result := string(data)
	for k, v := range values {
		result = strings.Replace(result, sep.start+k+sep.end, v, -1)
	}
	return []byte(result)
}

// ReplacingPlaceHoldersFromEnv replaces place holders in the string from
// Env vars
func ReplacingPlaceHoldersFromEnv(data []byte, sep Separator) ([]byte, error) {
	keys := listPlaceHolders(data, sep)
	values, err := getValuesFromEnv(keys)
	if err != nil {
		return nil, err
	}
	return ReplacingPlaceHolders(data, values, sep), nil
}

func getValuesFromEnv(keys []string) (map[string]string, error) {
	var notFound []string
	values := make(map[string]string, len(keys))
	for _, k := range keys {
		v, p := os.LookupEnv(k)
		if !p {
			notFound = append(notFound, k)
		} else {
			values[k] = v
		}
	}
	if len(notFound) > 0 {
		err := fmt.Errorf("Some values were not found: %+q", notFound)
		return values, err
	}
	return values, nil
}

// ReplacingPlaceHoldersInFilesFromEnv replaces placeholders in file from
// environment variables
func ReplacingPlaceHoldersInFilesFromEnv(files []string, sep Separator) error {
	keys := ListPlaceHoldersInFiles(files, sep)
	values, err := getValuesFromEnv(keys)
	if err != nil {
		return err
	}
	for _, f := range files {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		content = ReplacingPlaceHolders(content, values, sep)
		err = ioutil.WriteFile(f, content, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
