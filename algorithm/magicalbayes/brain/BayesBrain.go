package brain

import "fmt"

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