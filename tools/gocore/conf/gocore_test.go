package conf

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestGetGocoreConfig(t *testing.T) {
	c := GetGocoreConfig()
	_, err := yaml.Marshal(&c)
	if err != nil {
		t.Error(err)
	}
}
