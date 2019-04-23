package placeholder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/magiconair/properties"
)

// ValuesProvider define an interface exposing values
type ValuesProvider interface {
	getValue(key string) (string, bool)
}

// FileProvider implements ValuesProvider with values extracted from a json or yaml file
type FileProvider struct {
	values map[string]string
}

func (v FileProvider) getValue(key string) (string, bool) {
	val, ok := v.values[key]
	return val, ok
}

// NewFileProvider creates a new FileProvider from a json or yaml file
func NewFileProvider(input string) (*FileProvider, error) {
	file, err := ioutil.ReadFile(input)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(input, "yaml") {
		return newFileProviderYaml(file)
	}
	if strings.HasSuffix(input, "json") {
		return newFileProviderJSON(file)
	}
	if strings.HasSuffix(input, "properties") {
		return newFileProviderProperties(file)
	}
	return newFileProvider(file)
}

func newFileProviderJSON(content []byte) (*FileProvider, error) {
	var data interface{}
	ret := map[string]string{}
	err := json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	msg := data.(map[string]interface{})
	for k, v := range msg {
		ret[k] = v.(string)
	}
	return &FileProvider{values: ret}, nil
}

func newFileProviderProperties(content []byte) (*FileProvider, error) {
	p, err := properties.LoadString(string(content))
	if err != nil {
		return nil, err
	}
	return &FileProvider{values: p.Map()}, nil
}

func newFileProviderYaml(content []byte) (*FileProvider, error) {
	var data interface{}
	ret := map[string]string{}
	err := yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	msg := data.(map[interface{}]interface{})
	for k, v := range msg {
		key := k.(string)
		ret[key] = v.(string)
	}
	return &FileProvider{values: ret}, nil
}

func newFileProvider(content []byte) (*FileProvider, error) {
	// try with json
	provider, err := newFileProviderJSON(content)
	if err != nil {
		provider, err = newFileProviderYaml(content)
		if err != nil {
			return nil, err
		}
	}
	return provider, nil
}

// EnvProvider implements ValuesProvider with values extracted from environment
type EnvProvider struct{}

func (v EnvProvider) getValue(key string) (string, bool) {
	return os.LookupEnv(key)
}

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

func SeparatorFrom(start, end string) Separator {
	return Separator{
		start: regexp.QuoteMeta(start),
		end:   regexp.QuoteMeta(end),
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
func replacingPlaceHoldersFromValues(data []byte, values map[string]string, separator Separator) []byte {
	result := string(data)
	for k, v := range values {
		result = strings.Replace(result, separator.start+k+separator.end, v, -1)
	}
	return []byte(result)
}

// ReplacingPlaceHolders replaces place holders in the string from values provider
func ReplacingPlaceHolders(data []byte, separator Separator, provider ValuesProvider) ([]byte, error) {
	keys := listPlaceHolders(data, separator)
	values, err := getValues(keys, provider)
	if err != nil {
		return nil, err
	}
	return replacingPlaceHoldersFromValues(data, values, separator), nil
}

func getValues(keys []string, provider ValuesProvider) (map[string]string, error) {
	var notFound []string
	values := make(map[string]string, len(keys))
	for _, k := range keys {
		v, ok := provider.getValue(k)
		if !ok {
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

// ReplacingPlaceHoldersInFiles replaces placeholders in file from values provider
func ReplacingPlaceHoldersInFiles(files []string, separator Separator, provider ValuesProvider) error {
	keys := ListPlaceHoldersInFiles(files, separator)
	values, err := getValues(keys, provider)
	if err != nil {
		return err
	}
	for _, f := range files {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		content = replacingPlaceHoldersFromValues(content, values, separator)
		err = ioutil.WriteFile(f, content, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
