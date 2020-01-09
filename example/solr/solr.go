package main

import (
	"fmt"

	"github.com/sunmi-OS/gocore/solr"
)

func main() {

	si, _ := solr.NewSolrInterface("http://192.168.3.143:8983/solr", "app_collect")
	query := solr.NewQuery()
	query.Q("*:*")
	s := si.Search(query)
	r, err := s.Result(nil)
	if err != nil {
		fmt.Println(err.Error)
	} else {
		fmt.Println(r.Results.Docs)
	}
}
