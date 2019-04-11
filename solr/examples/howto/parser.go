package main

import "github.com/vanng822/go-solr/solr"
import "fmt"


type TestResultParser struct {
	original_response map[string]interface{}
}

func (parser *TestResultParser) Parse(response *solr.SelectResponse) (*solr.SolrResult, error) {
	sr := &solr.SolrResult{}
	sr.Results = new(solr.Collection)
	sr.Status = response.Status

	if response.Status == 0 {
		parser.ParseResponse(response, sr)
		parser.ParseFacetCounts(response, sr)
		parser.ParseHighlighting(response, sr)
	} else {
		parser.ParseError(response, sr)
	}
	parser.original_response = response.Response
	
	return sr, nil
}

func (parser *TestResultParser) ParseError(response *solr.SelectResponse, sr *solr.SolrResult) {
	if error, ok := response.Response["error"]; ok {
		sr.Error = error.(map[string]interface{})
	}
}


func (parser *TestResultParser) ParseResponse(response *solr.SelectResponse, sr *solr.SolrResult) {
	if resp, ok := response.Response["response"].(map[string]interface{}); ok {
		sr.Results.NumFound = int(resp["numFound"].(float64))
		sr.Results.Start = int(resp["start"].(float64))
		if docs, ok := resp["docs"].([]interface{}); ok {
			sr.Results.Docs = make([]solr.Document, len(docs))
			// remove version
			for i, v := range docs {
				d := solr.Document{}
				for k, v := range v.(map[string]interface{}) {
					if k != "_version_" {
						d.Set(k, v)
					}
				}
				sr.Results.Docs[i] = d
			}
		}
	} else {
		panic(`Standard parser can only parse solr response with response object,
					ie response.response and response.response.docs.
					Please use other parser or implement your own parser`)
	}
}

func (parser *TestResultParser) ParseFacetCounts(response *solr.SelectResponse, sr *solr.SolrResult) {
	if facetCounts, ok := response.Response["facet_counts"]; ok {
		sr.FacetCounts = facetCounts.(map[string]interface{})
	}
}

func (parser *TestResultParser) ParseHighlighting(response *solr.SelectResponse, sr *solr.SolrResult) {
	if highlighting, ok := response.Response["highlighting"]; ok {
		sr.Highlighting = highlighting.(map[string]interface{})
	}
}


type InheritResultParser struct {
	solr.StandardResultParser
	original_response map[string]interface{}
}


func (parser *InheritResultParser) Parse(response *solr.SelectResponse) (*solr.SolrResult, error) {
	sr, err := parser.StandardResultParser.Parse(response)
	if err != nil {
		return nil, err
	}
	parser.original_response = response.Response
	return sr, nil
}

func main() {
	si, _ := solr.NewSolrInterface("http://localhost:8983/solr", "collection1")

	query := solr.NewQuery()
	query.Q("title:add sucess 1")
	query.Start(0)
	query.Rows(15)
	s := si.Search(query)
	
	parser := &TestResultParser{}
	r, err := s.Result(parser)
	if err != nil {
		fmt.Println("Error when querying solr:", err.Error())
		return
	}
	
	fmt.Println(r.Results.Docs)
	fmt.Println(parser.original_response)
	
	parser2 := &InheritResultParser{}
	r2, err := s.Result(parser2)
	
	fmt.Println(r2.Results.Docs)
	
	fmt.Println(parser2.original_response)
}