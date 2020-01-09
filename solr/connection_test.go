package solr

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestConnection(t *testing.T) {
	_, err1 := NewConnection("fakedomain.tld", "core0")
	if err1 == nil {
		t.Errorf("It should be an error since the url is not valid")
	}

	_, err2 := NewConnection("http://www.fakedomain.tld", "core0")
	if err2 != nil {
		t.Errorf("It should not be an error since the url is  valid but got '%s'", err2.Error())
	}

	_, err3 := NewConnection("http://www.fakedomain.tld/select/", "core0")
	if err3 != nil {
		t.Errorf("It should not be an error since the url is  valid but got '%s'", err3.Error())
	}
}

func TestConnectionResourceInvalidDomain(t *testing.T) {
	conn, err := NewConnection("http://www.fakedomain.tld/", "core0")
	_, err = conn.Resource("select", &url.Values{})
	expected := "Get http://www.fakedomain.tld/core0/select?wt=json: dial tcp"
	error_report := err.Error()
	if strings.HasPrefix(error_report, expected) == false {
		t.Errorf("expected '%s' but got '%s'", expected, err.Error())
	}
}

func TestConnectionUpdateInvalidDomain(t *testing.T) {
	conn, err := NewConnection("http://www.fakedomain.tld/", "core0")
	_, err = conn.Update(map[string]interface{}{}, nil)
	expected := "Post http://www.fakedomain.tld/core0/update/?wt=json: dial tcp"
	error_report := err.Error()
	if strings.HasPrefix(error_report, expected) == false {
		t.Errorf("expected '%s' but got '%s'", expected, err.Error())
	}
}

func TestBytes2JsonWrongJson(t *testing.T) {
	data := []byte(`<xml><x>y</x><yy>boo</yy></xml>`)
	d, err := bytes2json(&data)
	if err == nil {
		t.Errorf("It should a error when parsing non json format")
	}
	if d != nil {
		t.Errorf("It should a error when parsing non json format")
	}
}

func TestBytes2Json(t *testing.T) {
	data := []byte(`{"t":"s","two":2,"obj":{"c":"b","j":"F"},"a":[1,2,3]}`)
	d, _ := bytes2json(&data)
	if d["t"] != "s" {
		t.Errorf("t should have s as value")
	}

	if d["two"].(float64) != 2 {
		t.Errorf("two should have 2 as value")
	}

}

func PrintMapInterface(d map[string]interface{}) {
	for k, v := range d {
		switch vv := v.(type) {
		case string:
			fmt.Println(fmt.Sprintf("%s:%s", k, v))
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float", vv)
		case map[string]interface{}:
			fmt.Println(k, "type is map[string]interface{}")
			PrintMapInterface(vv)
		case []interface{}:
			fmt.Println(k, "type is []interface{}")
			for i, u := range vv {
				switch uu := u.(type) {
				case map[string]interface{}:
					PrintMapInterface(uu)
				default:
					fmt.Println(i, u)
				}
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle", vv)
		}
	}
}

func TestJson2Bytes(t *testing.T) {

	test_json := map[string]interface{}{
		"t":   "s",
		"two": 2,
		"obj": map[string]interface{}{"c": "b", "j": "F"},
		"a":   []interface{}{1, 2, 3},
	}

	b, err := json2bytes(test_json)
	if err != nil {
		fmt.Println(err)
	}
	d, _ := bytes2json(b)

	if d["t"] != "s" {
		t.Errorf("t should have s as value")
	}

	if d["two"].(float64) != 2 {
		t.Errorf("two should have 2 as value")
	}
}

func TestHasError(t *testing.T) {
	data := map[string]interface{}{
		"responseHeader": map[string]interface{}{
			"status": float64(400),
			"QTime":  float64(30),
			"params": map[string]interface{}{
				"indent": "true",
				"q":      "*:*",
				"wt":     "json"}},
		"error": map[string]interface{}{
			"msg":  "no field name specified in query and no default specified via 'df' param",
			"code": float64(400)}}

	if hasError(data) != true {
		t.Errorf("Should have an error")
	}

	data2 := map[string]interface{}{
		"responseHeader": map[string]interface{}{
			"status": float64(0),
			"QTime":  float64(30),
			"params": map[string]interface{}{
				"indent": "true",
				"q":      "*:*",
				"wt":     "json"}},
		"response": map[string]interface{}{
			"numFound": float64(1),
			"start":    float64(0),
			"docs": []map[string]interface{}{{
				"id":        "change.me",
				"title":     "change.me",
				"_version_": float64(14)}}}}

	if hasError(data2) != false {
		t.Errorf("Should not has an error")
	}
}

func TestSuccessStatus(t *testing.T) {
	data := map[string]interface{}{
		"responseHeader": map[string]interface{}{
			"status": float64(400),
			"QTime":  float64(30),
			"params": map[string]interface{}{
				"indent": "true",
				"q":      "*:*",
				"wt":     "json"}},
		"error": map[string]interface{}{
			"msg":  "no field name specified in query and no default specified via 'df' param",
			"code": float64(400)}}
	if successStatus(data) != false {
		t.Errorf("Status check should give false but got true")
	}

	data2 := map[string]interface{}{
		"error": map[string]interface{}{
			"msg":  "Must specify a Content-Type header with POST requests",
			"code": float64(415)}}

	if successStatus(data2) != false {
		t.Errorf("Status check should give false but got true")
	}

	data3 := map[string]interface{}{
		"responseHeader": map[string]interface{}{
			"status": float64(0),
			"QTime":  float64(30),
			"params": map[string]interface{}{
				"indent": "true",
				"q":      "*:*",
				"wt":     "json"}},
		"response": map[string]interface{}{
			"numFound": float64(1),
			"start":    float64(0),
			"docs": []map[string]interface{}{{
				"id":        "change.me",
				"title":     "change.me",
				"_version_": float64(14)}}}}

	if successStatus(data3) != true {
		t.Errorf("Status check should give true but got false")
	}
}
