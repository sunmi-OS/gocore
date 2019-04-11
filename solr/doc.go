// Solr client in Go, core admin, add docs, update, delete, search and more 
//
//
//    package main
//    import (
//    	"github.com/vanng822/go-solr/solr"
//    	"fmt"
//    )
//  
//    func main() {
//      si, _ := solr.NewSolrInterface("http://localhost:8983/solr", "collection1")
//      query := solr.NewQuery()
//      query.Q("*:*")
//      s := si.Search(query)
//      r, _ := s.Result(nil)
//      fmt.Println(r.Results.Docs)
//    }
package solr
