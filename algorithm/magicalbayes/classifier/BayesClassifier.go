package classifier

import "github.com/sunmi-OS/gocore/algorithm/magicalbayes/brain"

var VERSION = "1.0.0"
var λ = 1 //平滑因子
var K = 1e3

// defaultProb is the tiny non-zero probability that a word
// we have not seen before appears in the class.
const defaultProb = 0.00000000001


type BayesClassifier struct {
	Brain *brain.BayesBrain
}

//func (classifier *BayesClassifier) probabilityOf(feature string, typesFrequency map[string]int) float64 {
//	frequency := typesFrequency[feature]
//	spaceSize := getSampleSpaceSize(typesFrequency)
//	//laplace平滑校准
//	frequency += λ
//	spaceSize += λ * len(classifier.Brain.FeaturesFrequency) // + λ * J
//	return float64(frequency) / float64(spaceSize)
//}

func (classifier *BayesClassifier) probabilityOf(feature string, typesFrequency map[string]float64) float64 {
	frequency := typesFrequency[feature]
	spaceSize := getSampleSpaceSize(typesFrequency)
	prob := float64(frequency) / float64(spaceSize)
	if prob < defaultProb {
		return defaultProb
	}
	return prob
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
	return classifier.probabilityOf(feature,  classifier.Brain.TfIdfTempValues[category])
}

func getSampleSpaceSize(typeFrequency map[string]float64) float64 {
	spaceSize := float64(0)
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


func (classifier *BayesClassifier) MolecularProbabilityOf(features ...string)  []Classification{
	i := 0
	list := make([]Classification, len(classifier.Brain.CategoriesFrequency))
	for category := range classifier.Brain.CategoriesFrequency {
		P := classifier.probabilityOfCategory(category)
		//fmt.Println("P", P)
		//fmt.Println("len",len(features), features)


		for _, feature := range features {
			inCategory := classifier.probabilityOfFeatureInCategory(feature, category)
			P *= inCategory * K
		}
		list[i] = Classification {
			Probability: P,
			Category:    category,
			Features:    features,
		}
		i++
	}
	return list

}

func (classifier *BayesClassifier) ProbabilityOf(features ...string) []Classification {
	i := 0
	list := make([]Classification, len(classifier.Brain.CategoriesFrequency))
	for category := range classifier.Brain.CategoriesFrequency {
		probability := classifier.BayesProbabilityOf(category, features...)
		list[i] = Classification {
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
			P *= classifier.probabilityOfFeatureInCategory(feature, category) * K
		}
		if P > maxProbability {
			maxProbability = P
			mostProbablyCategory = category
		}
	}
	return mostProbablyCategory
}

