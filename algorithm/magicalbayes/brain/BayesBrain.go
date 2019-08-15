package brain

import (
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"path/filepath"
)

var VERSION = "v0.0.1"

type BayesBrain struct {
	FeaturesFrequency               map[string]int
	CategoriesFrequency             map[string]int
	FeaturesFrequencyInEachCategory map[string]map[string]int
}

func NewBayesBrain() *BayesBrain{
	 brain := new(BayesBrain)
	 brain.FeaturesFrequency = make(map[string]int)
	 brain.CategoriesFrequency = make(map[string]int)
	 brain.FeaturesFrequencyInEachCategory = make(map[string]map[string]int)
	 return brain
}

func learn(featuresFrequency map[string]int,  features []string) {
	for _, feature := range features {
		featuresFrequency[feature]++
	}
}

func (brain *BayesBrain) Learn(category string, features ...string) {
	learn(brain.FeaturesFrequency, features)
	brain.CategoriesFrequency[category]++
	if brain.FeaturesFrequencyInEachCategory[category] == nil{
		brain.FeaturesFrequencyInEachCategory[category] = make(map[string]int)
	}
	learn(brain.FeaturesFrequencyInEachCategory[category], features)
}

func (brain *BayesBrain)Show() {
	fmt.Println("~~~~~~~~~~~ Bayes Brain " + VERSION + " ~~~~~~~~~~~")
	fmt.Println(brain.FeaturesFrequency)
	fmt.Println(brain.CategoriesFrequency)
	fmt.Println(brain.FeaturesFrequencyInEachCategory)
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
}

func (brain *BayesBrain) Save(filename string) error {
	if filename[len(filename):] != "/" {
		filename += "/"
	}
	parentPath := filename + VERSION
	err := save(brain.CategoriesFrequency, parentPath+"/cf.json")
	if err != nil {
		return err
	}
	err = save(brain.FeaturesFrequency, parentPath+"/ff.json")
	if err != nil {
		return err
	}
	err = save(brain.FeaturesFrequencyInEachCategory, parentPath+"/ffiec.json")
	if err != nil {
		return err
	}
	return nil
}

func save(obj interface{}, filename string) error {
	jsonBytes, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(obj)
	if err != nil {
		return err
	}

	path := filepath.Dir(filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, jsonBytes, 0600)
}



func (brain *BayesBrain) Load(filename string) error {
	if filename[len(filename):] != "/" {
		filename += "/"
	}
	parentPath := filename + VERSION
	jsonBytes, err := load(parentPath + "/cf.json")
	if err != nil {
		return err
	}
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(jsonBytes, &brain.CategoriesFrequency)
	if err != nil {
		return err
	}

	jsonBytes, err = load(parentPath + "/ff.json")
	if err != nil {
		return err
	}
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(jsonBytes, &brain.FeaturesFrequency)
	if err != nil {
		return err
	}


	jsonBytes, err = load(parentPath + "/ffiec.json")
	if err != nil {
		return err
	}
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(jsonBytes, &brain.FeaturesFrequencyInEachCategory)
	if err != nil {
		return err
	}

	return nil
}

func load(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return content, nil

}