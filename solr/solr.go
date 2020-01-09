package solr

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
)

// Shortcut for map[string]interface{}
// Use where applicable
type M map[string]interface{}

type Document map[string]interface{}

// Has check if a key exist in document
func (d Document) Has(k string) bool {
	_, ok := d[k]
	return ok
}

// Get returns value of a key
func (d Document) Get(k string) interface{} {
	v, _ := d[k]
	return v
}

// Set add a key/value to document
func (d Document) Set(k string, v interface{}) {
	d[k] = v
}

type SolrResponse struct {
	Status   int
	Response map[string]interface{}
}

type SolrUpdateResponse struct {
	Success bool
	Result  map[string]interface{}
}

// Holding the search result
type FireworkCollection struct {
	Docs     *json.RawMessage
	Start    int
	NumFound int
}

// Parsed result for SearchHandler response, ie /select
type FireworkSolrResult struct {
	Status  int                // status quick access to status
	Results FireworkCollection `json:"response"` // results parsed documents, basically response object
	QTime   int
	Params  map[string]string `json:"params"`

	ResponseHeader map[string]interface{}
	FacetCounts    map[string]interface{}

	Highlighting   map[string]interface{}
	Error          map[string]interface{}
	Grouped        map[string]interface{} // grouped for grouping result if grouping Results will be empty
	Stats          map[string]interface{}
	MoreLikeThis   map[string]interface{} // MoreLikeThis using Search (select) Component
	NextCursorMark string                 `json:"nextCursorMark"`
}

// Holding the search result
type Collection struct {
	Docs     []Document
	Start    int
	NumFound int
}

// Parsed result for SearchHandler response, ie /select
type SolrResult struct {
	Status         int         // status quick access to status
	Results        *Collection // results parsed documents, basically response object
	QTime          int
	ResponseHeader map[string]interface{}
	Facets         map[string]interface{}
	JsonFacets     map[string]interface{}
	FacetCounts    map[string]interface{}
	Highlighting   map[string]interface{}
	Error          map[string]interface{}
	Grouped        map[string]interface{} // grouped for grouping result if grouping Results will be empty
	Stats          map[string]interface{}
	MoreLikeThis   map[string]interface{} // MoreLikeThis using Search (select) Component
	SpellCheck     map[string]interface{} // SpellCheck using SpellCheck (spell) Component
	NextCursorMark string
}

// Parsed result for MoreLikeThisHandler response, ie /mlt
type SolrMltResult struct {
	Status         int         // status quick access to status
	Results        *Collection // results parsed documents, basically response object
	Match          *Collection // Documents for match section
	ResponseHeader map[string]interface{}
	Error          map[string]interface{}
}

type SolrInterface struct {
	conn *Connection
}

// Return a new instance of SolrInterface
func NewSolrInterface(solrUrl, core string) (*SolrInterface, error) {
	c, err := NewConnection(solrUrl, core)
	if err != nil {
		return nil, err
	}
	return &SolrInterface{conn: c}, nil
}

// Set to new core, this is just wrapper to Connection.SetCore which mean
// it will affect all places that use this Connection instance
func (si *SolrInterface) SetCore(core string) {
	si.conn.SetCore(core)
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
// See http://golang.org/pkg/net/http/#Request.SetBasicAuth
func (si *SolrInterface) SetBasicAuth(username, password string) {
	si.conn.SetBasicAuth(username, password)
}

// Return a new instace of Search, q is optional and one can set it later
func (si *SolrInterface) Search(q *Query) *Search {
	return NewSearch(si.conn, q)
}

// makeAddChunks splits the documents into chunks. If chunk_size is less than one it will be default to 100
func makeAddChunks(docs []Document, chunk_size int) []map[string]interface{} {
	if chunk_size < 1 {
		chunk_size = 100
	}
	docs_len := len(docs)
	num_chunk := int(math.Ceil(float64(docs_len) / float64(chunk_size)))
	doc_counter := 0
	chunks := make([]map[string]interface{}, num_chunk)
	for i := 0; i < num_chunk; i++ {
		add := make([]Document, 0, chunk_size)
		for j := 0; j < chunk_size; j++ {
			if doc_counter >= docs_len {
				break
			}
			add = append(add, docs[doc_counter])
			doc_counter++
		}
		chunks[i] = M{"add": add}
	}
	return chunks
}

// Add will insert documents in batch of chunk_size. success is false as long as one chunk failed.
// The result in SolrUpdateResponse is summery of response from all chunks
// with key chunk_%d
func (si *SolrInterface) Add(docs []Document, chunk_size int, params *url.Values) (*SolrUpdateResponse, error) {
	result := &SolrUpdateResponse{Success: true}
	responses := M{}
	chunks := makeAddChunks(docs, chunk_size)

	for i := 0; i < len(chunks); i++ {
		res, err := si.Update(chunks[i], params)
		if err != nil {
			return nil, err
		}
		result.Success = result.Success && res.Success
		responses[fmt.Sprintf("chunk_%d", i+1)] = M{
			"result":  res.Result,
			"success": res.Success,
			"total":   len(chunks[i]["add"].([]Document))}
	}
	result.Result = responses
	return result, nil
}

// Delete take data of type map and optional params which can use to specify addition parameters such as commit=true .
// Only one delete statement is supported, ie data can be { "id":"ID" } .
// If you want to delete more docs use { "query":"QUERY" } .
// Extra params can specify in params or in data such as { "query":"QUERY", "commitWithin":"500" }
func (si *SolrInterface) Delete(data map[string]interface{}, params *url.Values) (*SolrUpdateResponse, error) {
	message := M{"delete": data}
	return si.Update(message, params)
}

// DeleteAll will remove all documents and commit
func (si *SolrInterface) DeleteAll() (*SolrUpdateResponse, error) {
	params := &url.Values{}
	params.Add("commit", "true")
	return si.Delete(M{"query": "*:*"}, params)
}

// Update take data of type interface{} and optional params which can use to specify addition parameters such as commit=true
func (si *SolrInterface) Update(data interface{}, params *url.Values) (*SolrUpdateResponse, error) {
	if si.conn == nil {
		return nil, fmt.Errorf("No connection found for making request to solr")
	}
	return si.conn.Update(data, params)
}

// Commit the changes since the last commit
func (si *SolrInterface) Commit() (*SolrUpdateResponse, error) {
	params := &url.Values{}
	params.Add("commit", "true")
	return si.Update(M{}, params)
}

func (si *SolrInterface) Optimize(params *url.Values) (*SolrUpdateResponse, error) {
	if params == nil {
		params = &url.Values{}
	}
	params.Set("optimize", "true")
	return si.Update(M{}, params)
}

// Rollback rollbacks all add/deletes made to the index since the last commit.
// This should use with caution.
// See https://wiki.apache.org/solr/UpdateXmlMessages#A.22rollback.22
func (si *SolrInterface) Rollback() (*SolrUpdateResponse, error) {
	return si.Update(M{"rollback": M{}}, nil)
}

// Return new instance of CoreAdmin with provided solrUrl and basic auth
func (si *SolrInterface) CoreAdmin() (*CoreAdmin, error) {
	ca, err := NewCoreAdmin(si.conn.url.String())
	if err != nil {
		return nil, err
	}
	ca.SetBasicAuth(si.conn.username, si.conn.password)
	return ca, nil
}

// Return new instance of Schema with provided solrUrl and basic auth
func (si *SolrInterface) Schema() (*Schema, error) {
	s, err := NewSchema(si.conn.url.String(), si.conn.core)
	if err != nil {
		return nil, err
	}
	s.SetBasicAuth(si.conn.username, si.conn.password)
	return s, nil
}

// Return 'status' and QTime from solr, if everything is fine status should have value 'OK'
// QTime will have value -1 if can not determine
func (si *SolrInterface) Ping() (status string, qtime int, err error) {
	r, err := HTTPGet(fmt.Sprintf("%s/%s/admin/ping?wt=json", si.conn.url.String(), si.conn.core), nil, si.conn.username, si.conn.password)
	if err != nil {
		return "", -1, err
	}

	resp, err := bytes2json(&r)
	if err != nil {
		return "", -1, err
	}
	status, ok := resp["status"].(string)
	if ok == false {
		return "", -1, fmt.Errorf("Unexpected response returned")
	}
	if QTime, ok := resp["responseHeader"].(map[string]interface{})["QTime"].(float64); ok {
		qtime = int(QTime)
	} else {
		qtime = -1
	}
	return status, qtime, nil
}
