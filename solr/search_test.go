package solr

import "testing"


func TestSolrSearchMultipleValueSearchQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	q.AddParam("testing", "testing 2")
	res := s.QueryString()
	expected := "testing=test&testing=testing+2"
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
	}
}

func TestSolrSearchSetQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	expected := "testing=test"
	res := s.QueryString()
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
	}
	q2 := NewQuery()
	q2.AddParam("testing", "test2")
	s.SetQuery(q2)
	
	expected = "testing=test2"
	res = s.QueryString()
	
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
	}
}

func TestSolrSearchDebugQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	s.Debug = "true"
	res := s.QueryString()
	expected := "debug=true&indent=true&testing=test"
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
	}
}

func TestSolrSearchWithoutConnection(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)

	resp, err := s.Result(&StandardResultParser{})
	if resp != nil {
		t.Errorf("resp expected to be nil due to no connection is set")
	}
	if err == nil {
		t.Errorf("err expected to be not empty due to no connection is set")
	}
	expectedErrorMessage := "No connection found for making request to solr"

	if err.Error() != expectedErrorMessage {
		t.Errorf("The error message expecte to be '%s' but got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestSearchNoQuerySet(t *testing.T) {
	s := NewSearch(&Connection{}, nil)
	expected := ""
	if s.QueryString() != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, s.QueryString())
	}
}
