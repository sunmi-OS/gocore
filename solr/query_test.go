package solr

import "testing"

func TestSolrQueryAddParam(t *testing.T) {

	q := NewQuery()
	q.AddParam("qf", "some qf")

	if q.String() != "qf=some+qf" {
		t.Errorf("Expected to be: 'some qf'")
	}
}

func TestSolrQuerySetParam(t *testing.T) {

	q := NewQuery()
	q.SetParam("qf", "some qf")

	if q.String() != "qf=some+qf" {
		t.Errorf("Expected to be: 'some qf'")
	}
}

func TestSolrQueryGetParam(t *testing.T) {

	q := NewQuery()
	q.SetParam("qf", "some qf")

	if q.GetParam("qf") != "some qf" {
		t.Errorf("Expected to be: 'some qf'")
	}
}

func TestSolrQueryStart(t *testing.T) {

	q := NewQuery()
	q.Start(100)

	if q.String() != "start=100" {
		t.Errorf("Expected 'start=100'")
	}
}

func TestSolrSearchMultipleValueQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	q.AddParam("testing", "testing 2")
	res := q.String()
	expected := "testing=test&testing=testing+2"
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
	}
}

func TestSolrQueryRemoveParam(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	q.AddParam("testing2", "testing 2")
	// random order in for loop of range on map
	res := q.String()
	if res != "testing=test&testing2=testing+2" {
		t.Errorf("Expected to be: 'testing=test&testing2=testing+2' but got %s", res)
	}
	q.RemoveParam("testing2")
	if q.String() != "testing=test" {
		t.Errorf("Expected to be: 'testing=test'")
	}
}

func TestQueryQ(t *testing.T) {
	q := NewQuery()
	q.Q("id:100")
	expected := "q=id%3A100"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQuerySort(t *testing.T) {
	q := NewQuery()
	q.Sort("geodist() desc")
	expected := "sort=geodist%28%29+desc"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryFilterQuery(t *testing.T) {
	q := NewQuery()
	q.FilterQuery("popularity:[10 TO *]")
	expected := "fq=popularity%3A%5B10+TO+%2A%5D"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryFieldList(t *testing.T) {
	q := NewQuery()
	q.FieldList("id,name,decsription")
	expected := "fl=id%2Cname%2Cdecsription"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryGeofilt(t *testing.T) {
	q := NewQuery()
	q.Geofilt(45.15, -93.85, "store", 5)
	expected := "fq=%7B%21geofilt+pt%3D45.15%2C-93.85+sfield%3Dstore+d%3D5%7D"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryDefType(t *testing.T) {
	q := NewQuery()
	q.DefType("func")
	expected := "defType=func"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryBoostFunctions(t *testing.T) {
	q := NewQuery()
	q.BoostFunctions("recip(rord(myfield),1,2,3)")
	expected := "bf=recip%28rord%28myfield%29%2C1%2C2%2C3%29"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryBoostQuery(t *testing.T) {
	q := NewQuery()
	q.BoostQuery("cat:electronics^5.0")
	expected := "bq=cat%3Aelectronics%5E5.0"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryQueryField(t *testing.T) {
	q := NewQuery()
	q.QueryFields("features^20.0+text^0.3")
	expected := "qf=features%5E20.0%2Btext%5E0.3"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestSolrQueryAddFacet(t *testing.T) {
	q := NewQuery()
	q.AddFacet("color")
	q.AddFacet("size")
	expected := "facet=true&facet.field=color&facet.field=size"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestSolrQueryAddFacetQuery(t *testing.T) {
	q := NewQuery()
	q.AddFacetQuery("price:[* TO 500]")
	q.AddFacetQuery("price:[501 TO *]")
	expected := "facet=true&facet.query=price%3A%5B%2A+TO+500%5D&facet.query=price%3A%5B501+TO+%2A%5D"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestSolrQuerySetFacetMinCount(t *testing.T) {
	q := NewQuery()
	q.SetFacetMinCount(10)
	expected := "facet.mincount=10"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestSolrQueryAddFacetPivot(t *testing.T) {
	q := NewQuery()
	q.AddFacetPivot("color")
	q.AddFacetPivot("size")
	expected := "facet.pivot=color&facet.pivot=size"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestSolrQuerySetFacetPivotMinCount(t *testing.T) {
	q := NewQuery()
	q.SetFacetPivotMinCount(10)
	expected := "facet.pivot.mincount=10"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestSolrQueryAddJsonFacet(t *testing.T) {
	q := NewQuery()
	q.AddJsonFacet("{categories:{type:terms,field:cat,sort:'x desc',facet:{x:'avg(price)',y:'sum(price)'}}")
	expected := "json.facet=%7Bcategories%3A%7Btype%3Aterms%2Cfield%3Acat%2Csort%3A%27x+desc%27%2Cfacet%3A%7Bx%3A%27avg%28price%29%27%2Cy%3A%27sum%28price%29%27%7D%7D"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}
