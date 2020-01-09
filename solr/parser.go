package solr

import (
	"encoding/json"
	"fmt"
)

func bytes2json(data *[]byte) (map[string]interface{}, error) {
	var jsonData interface{}

	err := json.Unmarshal(*data, &jsonData)

	if err != nil {
		return nil, err
	}

	return jsonData.(map[string]interface{}), nil
}

// ResultParser is interface for parsing result from response.
// The idea here is that application have possibility to parse.
// Or defined own parser with internal data structure to suite
// application's need
type ResultParser interface {
	Parse(resp *[]byte) (*SolrResult, error)
}

type FireworkResultParser struct {
}

func (parser *FireworkResultParser) Parse(resp *[]byte) (FireworkSolrResult, error) {
	var res FireworkSolrResult
	err := json.Unmarshal(*resp, &res)
	return res, err
}

type ExtensiveResultParser struct {
}

func (parser *ExtensiveResultParser) Parse(resp_ *[]byte) (*SolrResult, error) {
	sr := &SolrResult{}
	jsonbuf, err := bytes2json(resp_)
	if err != nil {
		return sr, err
	}
	response := new(SolrResponse)
	response.Response = jsonbuf
	response.Status = int(jsonbuf["responseHeader"].(map[string]interface{})["status"].(float64))

	sr.Results = new(Collection)
	sr.Status = response.Status
	if nextCursorMark, ok := jsonbuf["nextCursorMark"]; ok {
		sr.NextCursorMark = fmt.Sprintf("%s", nextCursorMark)
	}

	parser.ParseResponseHeader(response, sr)

	if 0 != response.Status {
		parser.ParseError(response, sr)
		return sr, nil
	}

	err = parser.ParseResponse(response, sr)

	if err != nil {
		return nil, err
	}

	parser.ParseFacets(response, sr)
	parser.ParseJsonFacets(response, sr)

	return sr, nil
}

func (parser *ExtensiveResultParser) ParseResponseHeader(response *SolrResponse, sr *SolrResult) {
	if responseHeader, ok := response.Response["responseHeader"].(map[string]interface{}); ok {
		sr.ResponseHeader = responseHeader
	}
}

func (parser *ExtensiveResultParser) ParseError(response *SolrResponse, sr *SolrResult) {
	if err, ok := response.Response["error"].(map[string]interface{}); ok {
		sr.Error = err
	}
}

// ParseJsonFacets will assign facets and build sr.jsonfacets if there is a facet_counts
func (parser *ExtensiveResultParser) ParseFacets(response *SolrResponse, sr *SolrResult) {
	if fc, ok := response.Response["facet_counts"].(map[string]interface{}); ok {
		sr.FacetCounts = fc
		if f, ok := fc["facet_fields"].(map[string]interface{}); ok {
			sr.Facets = f
		}
	}
}

// ParseJsonFacets will assign facets and build sr.jsonfacets if there is a facets
func (parser *ExtensiveResultParser) ParseJsonFacets(response *SolrResponse, sr *SolrResult) {
	if jf, ok := response.Response["facets"].(map[string]interface{}); ok {
		sr.JsonFacets = jf
	}
}

// ParseSolrResponse will assign result and build sr.docs if there is a response.
// If there is no response or grouped property in response it will return error
func (parser *ExtensiveResultParser) ParseResponse(response *SolrResponse, sr *SolrResult) (err error) {
	if resp, ok := response.Response["response"].(map[string]interface{}); ok {
		ParseDocResponse(resp, sr.Results)
	} else {
		err = fmt.Errorf(`Extensive parser can only parse solr response with response object,
					ie response.response and response.response.docs. Or grouped response
					Please use other parser or implement your own parser`)
	}

	return err
}

type StandardResultParser struct {
}

func (parser *StandardResultParser) Parse(resp_ *[]byte) (*SolrResult, error) {

	sr := &SolrResult{}
	jsonbuf, err := bytes2json(resp_)
	if err != nil {
		return sr, err
	}
	response := new(SolrResponse)
	response.Response = jsonbuf
	response.Status = int(jsonbuf["responseHeader"].(map[string]interface{})["status"].(float64))

	sr.Results = new(Collection)
	sr.Status = response.Status
	if jsonbuf["nextCursorMark"] != nil {
		sr.NextCursorMark = fmt.Sprintf("%s", jsonbuf["nextCursorMark"])
	}

	parser.ParseResponseHeader(response, sr)

	if response.Status == 0 {
		err := parser.ParseResponse(response, sr)
		if err != nil {
			return nil, err
		}
		parser.ParseFacetCounts(response, sr)
		parser.ParseHighlighting(response, sr)
		parser.ParseStats(response, sr)
		parser.ParseMoreLikeThis(response, sr)
		parser.ParseSpellCheck(response, sr)
	} else {
		parser.ParseError(response, sr)
	}

	return sr, nil
}

func (parser *StandardResultParser) ParseResponseHeader(response *SolrResponse, sr *SolrResult) {
	if responseHeader, ok := response.Response["responseHeader"].(map[string]interface{}); ok {
		sr.ResponseHeader = responseHeader
	}
}

func (parser *StandardResultParser) ParseError(response *SolrResponse, sr *SolrResult) {
	if err, ok := response.Response["error"].(map[string]interface{}); ok {
		sr.Error = err
	}
}

func ParseDocResponse(docResponse map[string]interface{}, collection *Collection) {
	collection.NumFound = int(docResponse["numFound"].(float64))
	collection.Start = int(docResponse["start"].(float64))
	if docs, ok := docResponse["docs"].([]interface{}); ok {
		collection.Docs = make([]Document, len(docs))
		for i, v := range docs {
			collection.Docs[i] = Document(v.(map[string]interface{}))
		}
	}
}

// ParseSolrResponse will assign result and build sr.docs if there is a response.
// If there is no response or grouped property in response it will return error
func (parser *StandardResultParser) ParseResponse(response *SolrResponse, sr *SolrResult) (err error) {
	if resp, ok := response.Response["response"].(map[string]interface{}); ok {
		ParseDocResponse(resp, sr.Results)
	} else if grouped, ok := response.Response["grouped"].(map[string]interface{}); ok {
		sr.Grouped = grouped
	} else {
		err = fmt.Errorf(`Standard parser can only parse solr response with response object,
					ie response.response and response.response.docs. Or grouped response
					Please use other parser or implement your own parser`)
	}

	return err
}

// ParseFacetCounts will assign facet_counts to sr if there is one.
// No modification done here
func (parser *StandardResultParser) ParseFacetCounts(response *SolrResponse, sr *SolrResult) {
	if facetCounts, ok := response.Response["facet_counts"].(map[string]interface{}); ok {
		sr.FacetCounts = facetCounts
	}
}

// ParseHighlighting will assign highlighting to sr if there is one.
// No modification done here
func (parser *StandardResultParser) ParseHighlighting(response *SolrResponse, sr *SolrResult) {
	if highlighting, ok := response.Response["highlighting"].(map[string]interface{}); ok {
		sr.Highlighting = highlighting
	}
}

// Parse stats if there is  in response
func (parser *StandardResultParser) ParseStats(response *SolrResponse, sr *SolrResult) {
	if stats, ok := response.Response["stats"].(map[string]interface{}); ok {
		sr.Stats = stats
	}
}

// Parse moreLikeThis if there is in response
func (parser *StandardResultParser) ParseMoreLikeThis(response *SolrResponse, sr *SolrResult) {
	if moreLikeThis, ok := response.Response["moreLikeThis"].(map[string]interface{}); ok {
		sr.MoreLikeThis = moreLikeThis
	}
}

// Parse moreLikeThis if there is in response
func (parser *StandardResultParser) ParseSpellCheck(response *SolrResponse, sr *SolrResult) {
	if spellCheck, ok := response.Response["spellcheck"].(map[string]interface{}); ok {
		sr.SpellCheck = spellCheck
	}
}

type MltResultParser interface {
	Parse(*[]byte) (*SolrMltResult, error)
}

type MoreLikeThisParser struct {
}

func (parser *MoreLikeThisParser) Parse(resp_ *[]byte) (*SolrMltResult, error) {
	jsonbuf, err := bytes2json(resp_)
	sr := &SolrMltResult{}
	if err != nil {
		return sr, err
	}
	var resp = new(SolrResponse)
	resp.Response = jsonbuf
	resp.Status = int(jsonbuf["responseHeader"].(map[string]interface{})["status"].(float64))

	sr.Results = new(Collection)
	sr.Match = new(Collection)
	sr.Status = resp.Status

	if responseHeader, ok := resp.Response["responseHeader"].(map[string]interface{}); ok {
		sr.ResponseHeader = responseHeader
	}

	if resp.Status == 0 {
		if resp, ok := resp.Response["response"].(map[string]interface{}); ok {
			ParseDocResponse(resp, sr.Results)
		}
		if match, ok := resp.Response["match"].(map[string]interface{}); ok {
			ParseDocResponse(match, sr.Match)
		}
	} else {
		if err, ok := resp.Response["error"].(map[string]interface{}); ok {
			sr.Error = err
		}
	}
	return sr, nil
}
