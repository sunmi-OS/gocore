go-solr
=======


[![Build Status](https://travis-ci.org/vanng822/go-solr.svg?branch=master)](https://travis-ci.org/vanng822/go-solr)
[![GoDoc](https://godoc.org/github.com/vanng822/go-solr/solr?status.svg)](https://godoc.org/github.com/vanng822/go-solr/solr)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/vanng822/go-solr/solr) [![](http://gocover.io/_badge/github.com/vanng822/go-solr/solr)](http://gocover.io/github.com/vanng822/go-solr/solr)

Solr v4, required v4.4 if you want use all supported features.

Json only

No schema checking

Please go to http://wiki.apache.org/solr/ for how to write solr query.

## Features

Search, Add, Update, Delete, Commit, Rollback, Optimize

Core admin, Schema REST API


## Install

go get github.com/vanng822/go-solr/solr

## Usage

    package main
    import (
    	"github.com/vanng822/go-solr/solr"
    	"fmt"
    )
  
    func main() {
      si, _ := solr.NewSolrInterface("http://localhost:8983/solr", "collection1")
      query := solr.NewQuery()
      query.Q("*:*")
      s := si.Search(query)
      r, _ := s.Result(nil)
      fmt.Println(r.Results.Docs)
    }
    
## Developers

	export MOCK_LOGGING=1

for the mock logging

	unset MOCK_LOGGING

to remove this environment variable
	
## License
MIT
