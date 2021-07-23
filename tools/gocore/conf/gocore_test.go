package conf

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestGetGocoreConfig(t *testing.T) {
	c := GetGocoreConfig()
	s, err := yaml.Marshal(&c)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(s))
}
