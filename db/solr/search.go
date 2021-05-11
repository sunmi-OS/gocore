package solr

import (
	"fmt"
	"net/url"
)

type Search struct {
	query *Query
	conn  *Connection
	Debug string
}

// NewSearch takes c and q as optional
func NewSearch(c *Connection, q *Query) *Search {
	s := new(Search)
	if q != nil {
		s.SetQuery(q)
	}

	if c != nil {
		s.conn = c
	}
	return s
}

// SetQuery will replace old query with new query q
func (s *Search) SetQuery(q *Query) {
	s.query = q
}

// Return query params including debug and indent if Debug is set
func (s *Search) QueryParams() *url.Values {

	if s.query == nil {
		s.query = NewQuery()
	}

	if s.Debug != "" {
		s.query.params.Set("debug", s.Debug)
		s.query.params.Set("indent", "true")
	}

	return s.query.params
}

// QueryString return a query string of all queries except wt=json
func (s *Search) QueryString() string {
	return s.QueryParams().Encode()
}

// Wrapper for connection.Resource which will add wt=json automatically
// One can use this to query to /solr/{CORE}/{RESOURCE} example /solr/collection1/select
// This can be useful when you use an search component that is not supported in this package
func (s *Search) Resource(resource string, params *url.Values) (*[]byte, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("No connection found for making request to solr")
	}
	return s.conn.Resource(resource, params)
}

// Result will create a StandardResultParser if no parser specified.
// parser must be an implementation of ResultParser interface
func (s *Search) Result(parser ResultParser) (*SolrResult, error) {
	resp, err := s.Resource("select", s.QueryParams())
	if err != nil {
		return nil, err
	}
	if parser == nil {
		parser = new(StandardResultParser)
	}
	return parser.Parse(resp)
}

// This method is for making query to MoreLikeThisHandler
// See http://wiki.apache.org/solr/MoreLikeThisHandler
func (s *Search) MoreLikeThis(parser MltResultParser) (*SolrMltResult, error) {
	resp, err := s.Resource("mlt", s.QueryParams())
	if err != nil {
		return nil, err
	}
	if parser == nil {
		parser = new(MoreLikeThisParser)
	}
	return parser.Parse(resp)
}

// This method is for making query to SpellCheckHandler
// See https://wiki.apache.org/solr/SpellCheckComponent
func (s *Search) SpellCheck(parser ResultParser) (*SolrResult, error) {
	resp, err := s.Resource("spell", s.QueryParams())
	if err != nil {
		return nil, err
	}
	if parser == nil {
		parser = new(StandardResultParser)
	}
	return parser.Parse(resp)
}
