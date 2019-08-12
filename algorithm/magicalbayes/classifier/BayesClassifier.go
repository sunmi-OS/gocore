package classifier

import "gocore/algorithm/magicalbayes/brain"

var VERSION = "1.0.0"
var λ = 1 //平滑因子

type BayesClassifier struct {
	Brain * brain.BayesBrain
}

func (classifier *BayesClassifier) probabilityOf(feature string, typesFrequency map[string]int) float64 {
	frequency := typesFrequency[feature]
	spaceSize := getSampleSpaceSize(typesFrequency)
	//laplace平滑校准
	frequency += λ
	spaceSize += λ * len(classifier.Brain.FeaturesFrequency) // + λ * J
	return float64(frequency) / float64(spaceSize)
}

//P(feature)
func (classifier *BayesClassifier) probabilityOfFeature(feature string) float64 {
	return classifier.probabilityOf(feature, classifier.Brain.FeaturesFrequency)
}

//P(category)
func (classifier *BayesClassifier) probabilityOfCategory(category string) float64 {
 	categoryFrequency := classifier.Brain.CategoriesFrequency[category]
 	spaceSize := getSampleSpaceSize(classifier.Brain.CategoriesFrequency)
	return  float64(categoryFrequency) / float64(spaceSize)
}

//P(feature|category)
func (classifier *BayesClassifier) probabilityOfFeatureInCategory( feature, category string) float64 {
	return classifier.probabilityOf(feature,  classifier.Brain.FeaturesFrequencyInEachCategory[category])
}

func getSampleSpaceSize(typeFrequency map[string]int) int {
 	spaceSize := 0
	for _, frequency := range typeFrequency {
		spaceSize += frequency
	}
	return spaceSize
}

//P(category|features...)
//朴素贝叶斯计算特征数据属于某个类别的概率
func (classifier *BayesClassifier) BayesProbabilityOf(category string, features ...string) float64 {
	P := classifier.probabilityOfCategory(category)
	for _, feature := range features {
		P *= classifier.probabilityOfFeatureInCategory(feature, category)
	}
	space := 1.0

	for _, feature := range features {
		space *= classifier.probabilityOfFeature(feature)
	}
	return float64(P) / float64(space)
}


func (classifier *BayesClassifier) ProbabilityOf(features ...string) []Classification {
	i := 0
	list := make([]Classification, len(classifier.Brain.CategoriesFrequency))
	for category := range classifier.Brain.CategoriesFrequency {
		probability := classifier.BayesProbabilityOf(category, features...)
		list[i] = Classification{
			Probability: probability,
			Category:    category,
			Features:    features,
		}
		i++
	}
	return list
}


func (classifier *BayesClassifier) Classify(features ...string) string {
	mostProbablyCategory := ""
	maxProbability := -142857.0
	for category := range classifier.Brain.CategoriesFrequency {
		P := classifier.probabilityOfCategory(category)
		for _, feature := range features {
			P *= classifier.probabilityOfFeatureInCategory(feature, category)
		}
		if P > maxProbability {
			maxProbability = P
			mostProbablyCategory = category
		}
	}
	return mostProbablyCategory
}

