package solr

import (
	"fmt"
	"net/url"
	"strings"
)

type Schema struct {
	url      *url.URL
	core     string
	username string
	password string
}

// NewSchema will parse solrUrl and return a schema object, solrUrl must be a absolute url or path
func NewSchema(solrUrl, core string) (*Schema, error) {
	u, err := url.ParseRequestURI(strings.TrimRight(solrUrl, "/"))
	if err != nil {
		return nil, err
	}

	return &Schema{url: u, core: core}, nil
}

// Set to a new core
func (s *Schema) SetCore(core string) {
	s.core = core
}

func (s *Schema) SetBasicAuth(username, password string) {
	s.username = username
	s.password = password
}

// See Get requests in https://wiki.apache.org/solr/SchemaRESTAPI for detail
func (s *Schema) Get(path string, params *url.Values) (*SolrResponse, error) {
	var (
		r   []byte
		err error
	)
	if params == nil {
		params = &url.Values{}
	}

	params.Set("wt", "json")

	if path != "" {
		path = fmt.Sprintf("/%s", strings.Trim(path, "/"))
	}

	if s.core != "" {
		r, err = HTTPGet(fmt.Sprintf("%s/%s/schema%s?%s", s.url.String(), s.core, path, params.Encode()), nil, s.username, s.password)
	} else {
		r, err = HTTPGet(fmt.Sprintf("%s/schema%s?%s", s.url.String(), path, params.Encode()), nil, s.username, s.password)
	}
	if err != nil {
		return nil, err
	}
	resp, err := bytes2json(&r)
	if err != nil {
		return nil, err
	}

	return &SolrResponse{Response: resp, Status: int(resp["responseHeader"].(map[string]interface{})["status"].(float64))}, nil
}

//  Return entire schema, require Solr4.3, see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) All() (*SolrResponse, error) {
	return s.Get("", nil)
}

// Require Solr4.3, see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) Uniquekey() (*SolrResponse, error) {
	return s.Get("uniquekey", nil)
}

// Require Solr4.3, see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) Version() (*SolrResponse, error) {
	return s.Get("version", nil)
}

// Return name of schema, require Solr4.3, see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) Name() (*SolrResponse, error) {
	return s.Get("name", nil)
}

// see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) Fields(fl string, includeDynamic, showDefaults bool) (*SolrResponse, error) {
	params := &url.Values{}
	if includeDynamic {
		params.Set("includeDynamic", "true")
	}
	if showDefaults {
		params.Set("showDefaults", "true")
	}
	if fl != "" {
		params.Set("fl", fl)
	}
	return s.Get("fields", params)
}

// see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) FieldsName(name string, includeDynamic, showDefaults bool) (*SolrResponse, error) {
	params := &url.Values{}
	if includeDynamic {
		params.Set("includeDynamic", "true")
	}
	if showDefaults {
		params.Set("showDefaults", "true")
	}
	return s.Get(fmt.Sprintf("fields/%s", name), params)
}

// see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) Fieldtypes(showDefaults bool) (*SolrResponse, error) {
	params := &url.Values{}
	if showDefaults {
		params.Set("showDefaults", "true")
	}
	return s.Get("fieldtypes", params)
}

// see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) FieldtypesName(name string, showDefaults bool) (*SolrResponse, error) {
	params := &url.Values{}
	if showDefaults {
		params.Set("showDefaults", "true")
	}
	return s.Get(fmt.Sprintf("fieldtypes/%s", name), params)
}


// see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) DynamicFields(fl string, showDefaults bool) (*SolrResponse, error) {
	params := &url.Values{}
	if showDefaults {
		params.Set("showDefaults", "true")
	}
	if fl != "" {
		params.Set("fl", fl)
	}
	return s.Get("dynamicfields", params)
}

// see https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) DynamicFieldsName(name string, showDefaults bool) (*SolrResponse, error) {
	params := &url.Values{}
	if showDefaults {
		params.Set("showDefaults", "true")
	}
	return s.Get(fmt.Sprintf("dynamicfields/%s", name), params)
}

// For modify schema, require Solr4.4, currently one can add fields and copy fields.
// Example: s.Post("fields", data) for adding new fields.
// See https://wiki.apache.org/solr/SchemaRESTAPI
func (s *Schema) Post(path string, data interface{}) (*SolrUpdateResponse, error) {
	var (
		r   []byte
		err error
	)
	b, err := json2bytes(data)
	if err != nil {
		return nil, err
	}
	
	if s.core != "" {
		r, err = HTTPPost(fmt.Sprintf("%s/%s/schema/%s?wt=json", s.url.String(), s.core, strings.Trim(path, "/")), b, [][]string{{"Content-Type", "application/json"}}, s.username, s.password)
	} else {
		r, err = HTTPPost(fmt.Sprintf("%s/schema/%s?wt=json", s.url.String(), strings.Trim(path, "/")), b, [][]string{{"Content-Type", "application/json"}}, s.username, s.password)
	}
	if err != nil {
		return nil, err
	}
	
	resp, err := bytes2json(&r)
	if err != nil {
		return nil, err
	}
	// check error in resp
	if !successStatus(resp) || hasError(resp) {
		return &SolrUpdateResponse{Success: false, Result: resp}, nil
	}
	
	return &SolrUpdateResponse{Success: true, Result: resp}, nil
}
