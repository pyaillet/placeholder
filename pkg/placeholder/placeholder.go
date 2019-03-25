package placeholder

import (
	"io/ioutil"
	"regexp"
	"sort"

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
	re := regexp.MustCompile(sep.start + "(.*?)" + sep.end)

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
