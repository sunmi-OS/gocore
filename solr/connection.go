package solr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var userAgent = fmt.Sprintf("Go-solr/%s (+https://github.com/vanng822/go-solr)", VERSION)

var transport = http.Transport{}

// Solr imposes a limit on the size of a URL send to it using GET requests. Thus
// this library will switch to use to POST requests as the user query's grow up.
// If you need, you can charge this value, but be aware of the URL limit in
// your Solr distribution.
var MaximumSolrUrlLengthSupported = 2083

// HTTPPost make a POST request to path which also includes domain, headers are optional
func HTTPPost(path string, data *[]byte, headers [][]string, username, password string) ([]byte, error) {
	var (
		req *http.Request
		err error
	)

	client := &http.Client{Transport: &transport}
	if data == nil {
		req, err = http.NewRequest("POST", path, nil)
	} else {
		req, err = http.NewRequest("POST", path, bytes.NewReader(*data))
	}

	if err != nil {
		return nil, err
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	if len(headers) > 0 {
		for i := range headers {
			req.Header.Add(headers[i][0], headers[i][1])
		}
	}
	return makeRequest(client, req)
}

// HTTPGet make a GET request to url, headers are optional
func HTTPGet(url string, headers [][]string, username, password string) ([]byte, error) {
	client := &http.Client{Transport: &transport}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	if len(headers) > 0 {
		for i := range headers {
			req.Header.Add(headers[i][0], headers[i][1])
		}
	}
	return makeRequest(client, req)
}

func makeRequest(client *http.Client, req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func json2bytes(data interface{}) (*[]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func hasError(response map[string]interface{}) bool {
	_, ok := response["error"]
	return ok
}

func successStatus(response map[string]interface{}) bool {
	responseHeader, ok := response["responseHeader"].(map[string]interface{})
	if !ok {
		return false
	}

	if status, ok := responseHeader["status"].(float64); ok {
		return 0 == int(status)
	}

	return false
}

type Connection struct {
	url      *url.URL
	core     string
	username string
	password string
	headers [][]string
}

// NewConnection will parse solrUrl and return a connection object, solrUrl must be a absolute url or path
func NewConnection(solrUrl, core string) (*Connection, error) {
	u, err := url.ParseRequestURI(strings.TrimRight(solrUrl, "/"))
	if err != nil {
		return nil, err
	}
	return &Connection{url: u, core: core}, nil
}

// Set to a new core
func (c *Connection) SetCore(core string) {
	c.core = core
}

func (c *Connection) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
}

func (c *Connection) AddHeader(key, value string) {
	header := []string{key, value}
	c.headers = append(c.headers, header)
}

func (c *Connection) Resource(source string, params *url.Values) (*[]byte, error) {
	params.Set("wt", "json")
	baseUrl := fmt.Sprintf("%s/%s/%s", c.url.String(), c.core, source)
	encodedParameters := params.Encode()
	var r []byte
	var err error
	if len(baseUrl) + len(encodedParameters) >= MaximumSolrUrlLengthSupported {
		data := []byte(encodedParameters)
		var headers [][]string
		copy(headers, c.headers)
		headers = append(headers, []string{"Content-Type", "application/x-www-form-urlencoded"})
		r, err = HTTPPost(baseUrl, &data, headers, c.username, c.password)
	} else {
		r, err = HTTPGet(fmt.Sprintf("%s?%s", baseUrl, encodedParameters), c.headers, c.username, c.password)
	}
	return &r, err

}

// Update take optional params which can use to specify addition parameters such as commit=true
func (c *Connection) Update(data interface{}, params *url.Values) (*SolrUpdateResponse, error) {

	b, err := json2bytes(data)

	if err != nil {
		return nil, err
	}

	if params == nil {
		params = &url.Values{}
	}

	params.Set("wt", "json")

	r, err := HTTPPost(fmt.Sprintf("%s/%s/update/?%s", c.url.String(), c.core, params.Encode()), b, [][]string{{"Content-Type", "application/json"}}, c.username, c.password)

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
