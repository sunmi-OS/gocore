package solr

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func getBaciAuth(req *http.Request) (pair []string) {
	authData, ok := req.Header["Authorization"]
	if ok == false {
		return pair
	}
	auth := strings.SplitN(authData[0], " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return pair
	}
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair = strings.SplitN(string(payload), ":", 2)
	return pair
}

func authenticate(username, password string) bool {
	return username == "test" && password == "test"
}

func authenticateChallenge(w http.ResponseWriter, challenge string) {
	w.Header().Set("WWW-Authenticate", challenge)
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func logPrintBasicAuth(req *http.Request) {
	pair := getBaciAuth(req)
	log.Printf("Basic auth: %v", pair)
}

func logRequest(req *http.Request) {
	if os.Getenv("MOCK_LOGGING") != "" {
		log.Printf("RequestURI: %s", req.RequestURI)
		logPrintBasicAuth(req)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err.Error())
		}
		log.Println(string(body))
		log.Println(req.Header)
	}
}

func mockSuccessSelect(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":0,
		    "QTime":1,
		    "params":{
		      "indent":"true",
		      "q":"*:*",
		      "wt":"json"}},
		  "response":{"numFound":1,"start":0,"docs":[
		      {
		        "id":"change.me",
		        "title":["change.me"],
		        "_version_":1474699756018073600}]
		  }}`)
}

func mockSuccessSpell(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":0,
		    "QTime":1,
		    "params":{
		      "indent":"true",
		      "q":"*:*",
		      "wt":"json"}},
		  "response":{"numFound":0,"start":0,"docs":[]},
		  "spellcheck": {
            "suggestions": [
                  "tets",
                  {
                              "numFound": 5,
                              "startOffset": 0,
                              "endOffset": 5,
                              "origFreq": 3,
                              "suggestion": [{"word":"test","freq":9}]
                  }]}}`)
}

func mockSuccessSelectFacet(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
				  "responseHeader":{
				    "status":0,
				    "QTime":10,
				    "params":{
				      "facet":"true",
				      "indent":"true",
				      "q":"*:*",
				      "facet.field":"id",
				      "wt":"json"}},
				  "response":{"numFound":4,"start":0,"docs":[
				      {
				        "id":"change.me",
				        "title":["change.me"],
				        "_version_":1474893319511212032},
				      {
				        "id":"change.me2",
				        "title":["change.me2"],
				        "_version_":1474893328448225280},
				      {
				        "id":"change.me3",
				        "title":["change.me3"],
				        "_version_":1474893336208736256},
				      {
				        "id":"change.me2",
				        "title":["change.me22"],
				        "_version_":1474893362047746048}]
				  },
				  "facet_counts":{
				    "facet_queries":{},
				    "facet_fields":{
				      "id":[
				        "change.me2",2,
				        "change.me",1,
				        "change.me3",1]},
				    "facet_dates":{},
				    "facet_ranges":{}}}`)
}

func mockSuccessSelectHighlight(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
					  "responseHeader":{
					    "status":0,
					    "QTime":0,
					    "params":{
					      "indent":"true",
					      "q":"*:*",
					      "wt":"json",
					      "hl":"true"}},
					  "response":{"numFound":4,"start":0,"docs":[
					      {
					        "id":"change.me",
					        "title":["change.me"],
					        "_version_":1474893319511212032},
					      {
					        "id":"change.me2",
					        "title":["change.me2"],
					        "_version_":1474893328448225280},
					      {
					        "id":"change.me3",
					        "title":["change.me3"],
					        "_version_":1474893336208736256},
					      {
					        "id":"change.me2",
					        "title":["change.me22"],
					        "_version_":1474893362047746048}]
					  },
					  "highlighting":{
					    "change.me":{},
					    "change.me2":{},
					    "change.me3":{},
					    "change.me2":{}}}`)
}

func mockFailSelect(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	w.WriteHeader(400)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":400,
		    "QTime":3,
		    "params":{
		      "indent":"true",
		      "q":"**",
		      "wt":"json"}},
		  "error":{
		    "msg":"no field name specified in query and no default specified via 'df' param",
		    "code":400}}`)
}

func mockSuccessStandaloneCommit(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":5}}`)
}

func mockSuccessAdd(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	if req.Method == "POST" {
		if req.Header.Get("Content-Type") != "application/json" {
			writeContentTypeError(w)
			return
		}
	}
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":5}}`)
}

func mockSuccessDelete(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":5}}`)
}

// For commands that no need of specific response
func mockSuccessCommand(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":5}}`)
}

func mockSuccessGrouped(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":0,
		    "QTime":2,
		    "params":{
		      "q":"*:*",
		      "group.field":"id",
		      "group":"true",
		      "wt":"json"}},
		  "grouped":{
		    "id":{
		      "matches":2,
		      "groups":[{
		          "groupValue":"test_id_100",
		          "doclist":{"numFound":1,"start":0,"docs":[
		              {
		                "id":"test_id_100",
		                "title":["add sucess 100"],
		                "_version_":1475623982992457728}]
		          }},
		        {
		          "groupValue":"test_id_101",
		          "doclist":{"numFound":1,"start":0,"docs":[
		              {
		                "id":"test_id_101",
		                "title":["add sucess 101"],
		                "_version_":1475623982995603456}]
		          }}]}}}`)
}

func mockSuccessStats(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":0,
		    "QTime":1,
		    "params":{
		      "q":"*:*",
		      "stats":"true",
		      "indent":"true",
		      "rows":"0",
		      "wt":"json",
		      "stats.field":"id"}},
		  "response":{"numFound":200,"start":0,"docs":[]
		  },
		  "stats":{
		    "stats_fields":{
		      "id":{
		        "min":"test_id_0",
		        "max":"test_id_99",
		        "count":200,
		        "missing":0,
		        "facets":{}}}}}`)
}

func mockSuccessStrangeGrouped(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":0,
		    "QTime":2,
		    "params":{
		      "q":"*:*",
		      "group.field":"id",
		      "group":"true",
		      "wt":"json"}}}`)
}

func mockSuccessXML(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `<response>
						<lst name="responseHeader">
						<int name="status">0</int>
						<int name="QTime">8</int>
						</lst>
						</response>`)
}

func mockCoreAdmin(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":50}}`)
}

func mockSchema(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	login := getBaciAuth(req)
	if len(login) != 2 || !authenticate(login[0], login[1]) {
		authenticateChallenge(w, "Mock")
		return
	}
	
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":3},"schema":{"name":"example","version":1.5,"uniqueKey":"id","fieldTypes":[{"name":"alphaOnlySort","class":"solr.TextField","omitNorms":true,"sortMissingLast":true,"analyzer":{"tokenizer":{"class":"solr.KeywordTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.TrimFilterFactory"},{"class":"solr.PatternReplaceFilterFactory","pattern":"([^a-z])","replace":"all","replacement":""}]}},{"name":"ancestor_path","class":"solr.TextField","indexAnalyzer":{"tokenizer":{"class":"solr.KeywordTokenizerFactory"}},"queryAnalyzer":{"tokenizer":{"class":"solr.PathHierarchyTokenizerFactory","delimiter":"/"}}},{"name":"binary","class":"solr.BinaryField"},{"name":"boolean","class":"solr.BoolField","sortMissingLast":true},{"name":"currency","class":"solr.CurrencyField","currencyConfig":"currency.xml","defaultCurrency":"USD","precisionStep":"8"},{"name":"date","class":"solr.TrieDateField","positionIncrementGap":"0","precisionStep":"0"},{"name":"descendent_path","class":"solr.TextField","indexAnalyzer":{"tokenizer":{"class":"solr.PathHierarchyTokenizerFactory","delimiter":"/"}},"queryAnalyzer":{"tokenizer":{"class":"solr.KeywordTokenizerFactory"}}},{"name":"double","class":"solr.TrieDoubleField","positionIncrementGap":"0","precisionStep":"0"},{"name":"float","class":"solr.TrieFloatField","positionIncrementGap":"0","precisionStep":"0"},{"name":"ignored","class":"solr.StrField","indexed":false,"stored":false,"multiValued":true},{"name":"int","class":"solr.TrieIntField","positionIncrementGap":"0","precisionStep":"0"},{"name":"location","class":"solr.LatLonType","subFieldSuffix":"_coordinate"},{"name":"location_rpt","class":"solr.SpatialRecursivePrefixTreeFieldType","geo":"true","maxDistErr":"0.000009","distErrPct":"0.025","units":"degrees"},{"name":"long","class":"solr.TrieLongField","positionIncrementGap":"0","precisionStep":"0"},{"name":"lowercase","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.KeywordTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"}]}},{"name":"payloads","class":"solr.TextField","indexed":true,"stored":false,"analyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"},"filters":[{"class":"solr.DelimitedPayloadTokenFilterFactory","encoder":"float"}]}},{"name":"pdate","class":"solr.DateField","sortMissingLast":true},{"name":"pdouble","class":"solr.DoubleField"},{"name":"pfloat","class":"solr.FloatField"},{"name":"phonetic","class":"solr.TextField","indexed":true,"stored":false,"analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.DoubleMetaphoneFilterFactory","inject":"false"}]}},{"name":"pint","class":"solr.IntField"},{"name":"plong","class":"solr.LongField"},{"name":"point","class":"solr.PointType","subFieldSuffix":"_d","dimension":"2"},{"name":"random","class":"solr.RandomSortField","indexed":true},{"name":"string","class":"solr.StrField","sortMissingLast":true},{"name":"tdate","class":"solr.TrieDateField","positionIncrementGap":"0","precisionStep":"6"},{"name":"tdouble","class":"solr.TrieDoubleField","positionIncrementGap":"0","precisionStep":"8"},{"name":"text_ar","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ar.txt","ignoreCase":"true"},{"class":"solr.ArabicNormalizationFilterFactory"},{"class":"solr.ArabicStemFilterFactory"}]}},{"name":"text_bg","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_bg.txt","ignoreCase":"true"},{"class":"solr.BulgarianStemFilterFactory"}]}},{"name":"text_ca","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.ElisionFilterFactory","articles":"lang/contractions_ca.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ca.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Catalan"}]}},{"name":"text_cjk","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.CJKWidthFilterFactory"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.CJKBigramFilterFactory"}]}},{"name":"text_ckb","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.SoraniNormalizationFilterFactory"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ckb.txt","ignoreCase":"true"},{"class":"solr.SoraniStemFilterFactory"}]}},{"name":"text_cz","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_cz.txt","ignoreCase":"true"},{"class":"solr.CzechStemFilterFactory"}]}},{"name":"text_da","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_da.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Danish"}]}},{"name":"text_de","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_de.txt","ignoreCase":"true"},{"class":"solr.GermanNormalizationFilterFactory"},{"class":"solr.GermanLightStemFilterFactory"}]}},{"name":"text_el","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.GreekLowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_el.txt","ignoreCase":"false"},{"class":"solr.GreekStemFilterFactory"}]}},{"name":"text_en","class":"solr.TextField","positionIncrementGap":"100","indexAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.EnglishPossessiveFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.PorterStemFilterFactory"}]},"queryAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.SynonymFilterFactory","expand":"true","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.EnglishPossessiveFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.PorterStemFilterFactory"}]}},{"name":"text_en_splitting","class":"solr.TextField","autoGeneratePhraseQueries":"true","positionIncrementGap":"100","indexAnalyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.WordDelimiterFilterFactory","catenateNumbers":"1","generateNumberParts":"1","splitOnCaseChange":"1","generateWordParts":"1","catenateAll":"0","catenateWords":"1"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.PorterStemFilterFactory"}]},"queryAnalyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"},"filters":[{"class":"solr.SynonymFilterFactory","expand":"true","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.WordDelimiterFilterFactory","catenateNumbers":"0","generateNumberParts":"1","splitOnCaseChange":"1","generateWordParts":"1","catenateAll":"0","catenateWords":"0"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.PorterStemFilterFactory"}]}},{"name":"text_en_splitting_tight","class":"solr.TextField","autoGeneratePhraseQueries":"true","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"},"filters":[{"class":"solr.SynonymFilterFactory","expand":"false","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.WordDelimiterFilterFactory","catenateNumbers":"1","generateNumberParts":"0","generateWordParts":"0","catenateAll":"0","catenateWords":"1"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.EnglishMinimalStemFilterFactory"},{"class":"solr.RemoveDuplicatesTokenFilterFactory"}]}},{"name":"text_es","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_es.txt","ignoreCase":"true"},{"class":"solr.SpanishLightStemFilterFactory"}]}},{"name":"text_eu","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_eu.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Basque"}]}},{"name":"text_fa","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"charFilters":[{"class":"solr.PersianCharFilterFactory"}],"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.ArabicNormalizationFilterFactory"},{"class":"solr.PersianNormalizationFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_fa.txt","ignoreCase":"true"}]}},{"name":"text_fi","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_fi.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Finnish"}]}},{"name":"text_fr","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.ElisionFilterFactory","articles":"lang/contractions_fr.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_fr.txt","ignoreCase":"true"},{"class":"solr.FrenchLightStemFilterFactory"}]}},{"name":"text_ga","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.ElisionFilterFactory","articles":"lang/contractions_ga.txt","ignoreCase":"true"},{"class":"solr.StopFilterFactory","words":"lang/hyphenations_ga.txt","ignoreCase":"true"},{"class":"solr.IrishLowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ga.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Irish"}]}},{"name":"text_general","class":"solr.TextField","positionIncrementGap":"100","indexAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"stopwords.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"}]},"queryAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"stopwords.txt","ignoreCase":"true"},{"class":"solr.SynonymFilterFactory","expand":"true","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.LowerCaseFilterFactory"}]}},{"name":"text_general_rev","class":"solr.TextField","positionIncrementGap":"100","indexAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"stopwords.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.ReversedWildcardFilterFactory","maxPosQuestion":"2","maxFractionAsterisk":"0.33","maxPosAsterisk":"3","withOriginal":"true"}]},"queryAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.SynonymFilterFactory","expand":"true","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.StopFilterFactory","words":"stopwords.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"}]}},{"name":"text_gl","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_gl.txt","ignoreCase":"true"},{"class":"solr.GalicianStemFilterFactory"}]}},{"name":"text_hi","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.IndicNormalizationFilterFactory"},{"class":"solr.HindiNormalizationFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_hi.txt","ignoreCase":"true"},{"class":"solr.HindiStemFilterFactory"}]}},{"name":"text_hu","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_hu.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Hungarian"}]}},{"name":"text_hy","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_hy.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Armenian"}]}},{"name":"text_id","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_id.txt","ignoreCase":"true"},{"class":"solr.IndonesianStemFilterFactory","stemDerivational":"true"}]}},{"name":"text_it","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.ElisionFilterFactory","articles":"lang/contractions_it.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_it.txt","ignoreCase":"true"},{"class":"solr.ItalianLightStemFilterFactory"}]}},{"name":"text_ja","class":"solr.TextField","autoGeneratePhraseQueries":"false","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.JapaneseTokenizerFactory","mode":"search"},"filters":[{"class":"solr.JapaneseBaseFormFilterFactory"},{"class":"solr.JapanesePartOfSpeechStopFilterFactory","tags":"lang/stoptags_ja.txt"},{"class":"solr.CJKWidthFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ja.txt","ignoreCase":"true"},{"class":"solr.JapaneseKatakanaStemFilterFactory","minimumLength":"4"},{"class":"solr.LowerCaseFilterFactory"}]}},{"name":"text_lv","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_lv.txt","ignoreCase":"true"},{"class":"solr.LatvianStemFilterFactory"}]}},{"name":"text_nl","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_nl.txt","ignoreCase":"true"},{"class":"solr.StemmerOverrideFilterFactory","dictionary":"lang/stemdict_nl.txt","ignoreCase":"false"},{"class":"solr.SnowballPorterFilterFactory","language":"Dutch"}]}},{"name":"text_no","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_no.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Norwegian"}]}},{"name":"text_pt","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_pt.txt","ignoreCase":"true"},{"class":"solr.PortugueseLightStemFilterFactory"}]}},{"name":"text_ro","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ro.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Romanian"}]}},{"name":"text_ru","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_ru.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Russian"}]}},{"name":"text_sv","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_sv.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Swedish"}]}},{"name":"text_th","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.ThaiWordFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_th.txt","ignoreCase":"true"}]}},{"name":"text_tr","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.TurkishLowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_tr.txt","ignoreCase":"false"},{"class":"solr.SnowballPorterFilterFactory","language":"Turkish"}]}},{"name":"text_ws","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"}}},{"name":"tfloat","class":"solr.TrieFloatField","positionIncrementGap":"0","precisionStep":"8"},{"name":"tint","class":"solr.TrieIntField","positionIncrementGap":"0","precisionStep":"8"},{"name":"tlong","class":"solr.TrieLongField","positionIncrementGap":"0","precisionStep":"8"}],"fields":[{"name":"_root_","type":"string","indexed":true,"stored":false},{"name":"_version_","type":"long","indexed":true,"stored":true},{"name":"author","type":"text_general","indexed":true,"stored":true},{"name":"cat","type":"string","multiValued":true,"indexed":true,"stored":true},{"name":"category","type":"text_general","indexed":true,"stored":true},{"name":"comments","type":"text_general","indexed":true,"stored":true},{"name":"content","type":"text_general","multiValued":true,"indexed":false,"stored":true},{"name":"content_type","type":"string","multiValued":true,"indexed":true,"stored":true},{"name":"description","type":"text_general","indexed":true,"stored":true},{"name":"features","type":"text_general","multiValued":true,"indexed":true,"stored":true},{"name":"id","type":"string","multiValued":false,"indexed":true,"required":true,"stored":true},{"name":"inStock","type":"boolean","indexed":true,"stored":true},{"name":"includes","type":"text_general","termPositions":true,"termVectors":true,"indexed":true,"termOffsets":true,"stored":true},{"name":"keywords","type":"text_general","indexed":true,"stored":true},{"name":"last_modified","type":"date","indexed":true,"stored":true},{"name":"links","type":"string","multiValued":true,"indexed":true,"stored":true},{"name":"manu","type":"text_general","omitNorms":true,"indexed":true,"stored":true},{"name":"manu_exact","type":"string","indexed":true,"stored":false},{"name":"name","type":"text_general","indexed":true,"stored":true},{"name":"payloads","type":"payloads","indexed":true,"stored":true},{"name":"popularity","type":"int","indexed":true,"stored":true},{"name":"price","type":"float","indexed":true,"stored":true},{"name":"resourcename","type":"text_general","indexed":true,"stored":true},{"name":"sku","type":"text_en_splitting_tight","omitNorms":true,"indexed":true,"stored":true},{"name":"store","type":"location","indexed":true,"stored":true},{"name":"subject","type":"text_general","indexed":true,"stored":true},{"name":"text","type":"text_general","multiValued":true,"indexed":true,"stored":false},{"name":"text_rev","type":"text_general_rev","multiValued":true,"indexed":true,"stored":false},{"name":"title","type":"text_general","multiValued":true,"indexed":true,"stored":true},{"name":"url","type":"text_general","indexed":true,"stored":true},{"name":"weight","type":"float","indexed":true,"stored":true}],"dynamicFields":[{"name":"*_coordinate","type":"tdouble","indexed":true,"stored":false},{"name":"ignored_*","type":"ignored","multiValued":true},{"name":"random_*","type":"random"},{"name":"attr_*","type":"text_general","multiValued":true,"indexed":true,"stored":true},{"name":"*_txt","type":"text_general","multiValued":true,"indexed":true,"stored":true},{"name":"*_dts","type":"date","multiValued":true,"indexed":true,"stored":true},{"name":"*_tdt","type":"tdate","indexed":true,"stored":true},{"name":"*_is","type":"int","multiValued":true,"indexed":true,"stored":true},{"name":"*_ss","type":"string","multiValued":true,"indexed":true,"stored":true},{"name":"*_ls","type":"long","multiValued":true,"indexed":true,"stored":true},{"name":"*_en","type":"text_en","multiValued":true,"indexed":true,"stored":true},{"name":"*_bs","type":"boolean","multiValued":true,"indexed":true,"stored":true},{"name":"*_fs","type":"float","multiValued":true,"indexed":true,"stored":true},{"name":"*_ds","type":"double","multiValued":true,"indexed":true,"stored":true},{"name":"*_dt","type":"date","indexed":true,"stored":true},{"name":"*_ti","type":"tint","indexed":true,"stored":true},{"name":"*_tl","type":"tlong","indexed":true,"stored":true},{"name":"*_tf","type":"tfloat","indexed":true,"stored":true},{"name":"*_td","type":"tdouble","indexed":true,"stored":true},{"name":"*_pi","type":"pint","indexed":true,"stored":true},{"name":"*_i","type":"int","indexed":true,"stored":true},{"name":"*_s","type":"string","indexed":true,"stored":true},{"name":"*_l","type":"long","indexed":true,"stored":true},{"name":"*_t","type":"text_general","indexed":true,"stored":true},{"name":"*_b","type":"boolean","indexed":true,"stored":true},{"name":"*_f","type":"float","indexed":true,"stored":true},{"name":"*_d","type":"double","indexed":true,"stored":true},{"name":"*_p","type":"location","indexed":true,"stored":true},{"name":"*_c","type":"currency","indexed":true,"stored":true}],"copyFields":[{"source":"author","dest":"text"},{"source":"cat","dest":"text"},{"source":"content","dest":"text"},{"source":"content_type","dest":"text"},{"source":"description","dest":"text"},{"source":"features","dest":"text"},{"source":"includes","dest":"text"},{"source":"keywords","dest":"text"},{"source":"manu","dest":"manu_exact"},{"source":"manu","dest":"text"},{"source":"name","dest":"text"},{"source":"resourcename","dest":"text"},{"source":"title","dest":"text"},{"source":"url","dest":"text"},{"source":"price","dest":"price_c"},{"source":"author","dest":"author_s"}]}}`)
}

func mockSchemaFields(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	if req.Method == "POST" {
		if req.Header.Get("Content-Type") != "application/json" {
			writeContentTypeError(w)
			return
		}
	}
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":1},"fields":[{"name":"_root_","type":"string","indexed":true,"stored":false},{"name":"_version_","type":"long","indexed":true,"stored":true},{"name":"author","type":"text_general","indexed":true,"stored":true},{"name":"cat","type":"string","multiValued":true,"indexed":true,"stored":true},{"name":"category","type":"text_general","indexed":true,"stored":true},{"name":"comments","type":"text_general","indexed":true,"stored":true},{"name":"content","type":"text_general","multiValued":true,"indexed":false,"stored":true},{"name":"content_type","type":"string","multiValued":true,"indexed":true,"stored":true},{"name":"description","type":"text_general","indexed":true,"stored":true},{"name":"features","type":"text_general","multiValued":true,"indexed":true,"stored":true},{"name":"id","type":"string","multiValued":false,"indexed":true,"required":true,"stored":true,"uniqueKey":true},{"name":"inStock","type":"boolean","indexed":true,"stored":true},{"name":"includes","type":"text_general","termPositions":true,"termVectors":true,"indexed":true,"termOffsets":true,"stored":true},{"name":"keywords","type":"text_general","indexed":true,"stored":true},{"name":"last_modified","type":"date","indexed":true,"stored":true},{"name":"links","type":"string","multiValued":true,"indexed":true,"stored":true},{"name":"manu","type":"text_general","omitNorms":true,"indexed":true,"stored":true},{"name":"manu_exact","type":"string","indexed":true,"stored":false},{"name":"name","type":"text_general","indexed":true,"stored":true},{"name":"payloads","type":"payloads","indexed":true,"stored":true},{"name":"popularity","type":"int","indexed":true,"stored":true},{"name":"price","type":"float","indexed":true,"stored":true},{"name":"resourcename","type":"text_general","indexed":true,"stored":true},{"name":"sku","type":"text_en_splitting_tight","omitNorms":true,"indexed":true,"stored":true},{"name":"store","type":"location","indexed":true,"stored":true},{"name":"subject","type":"text_general","indexed":true,"stored":true},{"name":"text","type":"text_general","multiValued":true,"indexed":true,"stored":false},{"name":"text_rev","type":"text_general_rev","multiValued":true,"indexed":true,"stored":false},{"name":"title","type":"text_general","multiValued":true,"indexed":true,"stored":true},{"name":"url","type":"text_general","indexed":true,"stored":true},{"name":"weight","type":"float","indexed":true,"stored":true}]}`)
}

func mockSchemaFieldsTitle(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":2},"field":{"name":"title","type":"text_general","multiValued":true,"indexed":true,"stored":true}}`)
}

func mockSchemaDynamicFields(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":2},"dynamicFields":[{"name":"*_coordinate","type":"tdouble","indexed":true,"stored":false},{"name":"ignored_*","type":"ignored","multiValued":true},{"name":"random_*","type":"random"},{"name":"attr_*","type":"text_general","multiValued":true,"indexed":true,"stored":true},{"name":"*_txt","type":"text_general","multiValued":true,"indexed":true,"stored":true},{"name":"*_dts","type":"date","multiValued":true,"indexed":true,"stored":true},{"name":"*_tdt","type":"tdate","indexed":true,"stored":true},{"name":"*_is","type":"int","multiValued":true,"indexed":true,"stored":true},{"name":"*_ss","type":"string","multiValued":true,"indexed":true,"stored":true},{"name":"*_ls","type":"long","multiValued":true,"indexed":true,"stored":true},{"name":"*_en","type":"text_en","multiValued":true,"indexed":true,"stored":true},{"name":"*_bs","type":"boolean","multiValued":true,"indexed":true,"stored":true},{"name":"*_fs","type":"float","multiValued":true,"indexed":true,"stored":true},{"name":"*_ds","type":"double","multiValued":true,"indexed":true,"stored":true},{"name":"*_dt","type":"date","indexed":true,"stored":true},{"name":"*_ti","type":"tint","indexed":true,"stored":true},{"name":"*_tl","type":"tlong","indexed":true,"stored":true},{"name":"*_tf","type":"tfloat","indexed":true,"stored":true},{"name":"*_td","type":"tdouble","indexed":true,"stored":true},{"name":"*_pi","type":"pint","indexed":true,"stored":true},{"name":"*_i","type":"int","indexed":true,"stored":true},{"name":"*_s","type":"string","indexed":true,"stored":true},{"name":"*_l","type":"long","indexed":true,"stored":true},{"name":"*_t","type":"text_general","indexed":true,"stored":true},{"name":"*_b","type":"boolean","indexed":true,"stored":true},{"name":"*_f","type":"float","indexed":true,"stored":true},{"name":"*_d","type":"double","indexed":true,"stored":true},{"name":"*_p","type":"location","indexed":true,"stored":true},{"name":"*_c","type":"currency","indexed":true,"stored":true}]}`)
}

func mockSchemaDynamicFieldsCoordinate(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":1},"dynamicField":{"name":"*_coordinate","type":"tdouble","indexed":true,"stored":false}}`)
}

func mockSchemaFieldTypes(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":3},"fieldTypes":[{"name":"alphaOnlySort","class":"solr.TextField","omitNorms":true,"sortMissingLast":true,"analyzer":{"tokenizer":{"class":"solr.KeywordTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.TrimFilterFactory"},{"class":"solr.PatternReplaceFilterFactory","pattern":"([^a-z])","replace":"all","replacement":""}]},"fields":[],"dynamicFields":[]},{"name":"ancestor_path","class":"solr.TextField","indexAnalyzer":{"tokenizer":{"class":"solr.KeywordTokenizerFactory"}},"queryAnalyzer":{"tokenizer":{"class":"solr.PathHierarchyTokenizerFactory","delimiter":"/"}},"fields":[],"dynamicFields":[]},{"name":"binary","class":"solr.BinaryField","fields":[],"dynamicFields":[]},{"name":"boolean","class":"solr.BoolField","sortMissingLast":true,"fields":["inStock"],"dynamicFields":["*_bs","*_b"]},{"name":"currency","class":"solr.CurrencyField","currencyConfig":"currency.xml","defaultCurrency":"USD","precisionStep":"8","fields":[],"dynamicFields":["*_c"]},{"name":"date","class":"solr.TrieDateField","positionIncrementGap":"0","precisionStep":"0","fields":["last_modified"],"dynamicFields":["*_dts","*_dt"]},{"name":"descendent_path","class":"solr.TextField","indexAnalyzer":{"tokenizer":{"class":"solr.PathHierarchyTokenizerFactory","delimiter":"/"}},"queryAnalyzer":{"tokenizer":{"class":"solr.KeywordTokenizerFactory"}},"fields":[],"dynamicFields":[]},{"name":"double","class":"solr.TrieDoubleField","positionIncrementGap":"0","precisionStep":"0","fields":[],"dynamicFields":["*_ds","*_d"]},{"name":"float","class":"solr.TrieFloatField","positionIncrementGap":"0","precisionStep":"0","fields":["price","weight"],"dynamicFields":["*_fs","*_f"]},{"name":"ignored","class":"solr.StrField","indexed":false,"stored":false,"multiValued":true,"fields":[],"dynamicFields":["ignored_*"]},{"name":"int","class":"solr.TrieIntField","positionIncrementGap":"0","precisionStep":"0","fields":["popularity"],"dynamicFields":["*_is","*_i"]},{"name":"location","class":"solr.LatLonType","subFieldSuffix":"_coordinate","fields":["store"],"dynamicFields":["*_p"]},{"name":"location_rpt","class":"solr.SpatialRecursivePrefixTreeFieldType","geo":"true","maxDistErr":"0.000009","distErrPct":"0.025","units":"degrees","fields":[],"dynamicFields":[]},{"name":"long","class":"solr.TrieLongField","positionIncrementGap":"0","precisionStep":"0","fields":["_version_"],"dynamicFields":["*_ls","*_l"]},{"name":"lowercase","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.KeywordTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"payloads","class":"solr.TextField","indexed":true,"stored":false,"analyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"},"filters":[{"class":"solr.DelimitedPayloadTokenFilterFactory","encoder":"float"}]},"fields":["payloads"],"dynamicFields":[]},{"name":"pdate","class":"solr.DateField","sortMissingLast":true,"fields":[],"dynamicFields":[]},{"name":"pdouble","class":"solr.DoubleField","fields":[],"dynamicFields":[]},{"name":"pfloat","class":"solr.FloatField","fields":[],"dynamicFields":[]},{"name":"phonetic","class":"solr.TextField","indexed":true,"stored":false,"analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.DoubleMetaphoneFilterFactory","inject":"false"}]},"fields":[],"dynamicFields":[]},{"name":"pint","class":"solr.IntField","fields":[],"dynamicFields":["*_pi"]},{"name":"plong","class":"solr.LongField","fields":[],"dynamicFields":[]},{"name":"point","class":"solr.PointType","subFieldSuffix":"_d","dimension":"2","fields":[],"dynamicFields":[]},{"name":"random","class":"solr.RandomSortField","indexed":true,"fields":[],"dynamicFields":["random_*"]},{"name":"string","class":"solr.StrField","sortMissingLast":true,"fields":["_root_","cat","content_type","id","links","manu_exact"],"dynamicFields":["*_ss","*_s"]},{"name":"tdate","class":"solr.TrieDateField","positionIncrementGap":"0","precisionStep":"6","fields":[],"dynamicFields":["*_tdt"]},{"name":"tdouble","class":"solr.TrieDoubleField","positionIncrementGap":"0","precisionStep":"8","fields":[],"dynamicFields":["*_coordinate","*_td"]},{"name":"text_ar","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ar.txt","ignoreCase":"true"},{"class":"solr.ArabicNormalizationFilterFactory"},{"class":"solr.ArabicStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_bg","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_bg.txt","ignoreCase":"true"},{"class":"solr.BulgarianStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_ca","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.ElisionFilterFactory","articles":"lang/contractions_ca.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ca.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Catalan"}]},"fields":[],"dynamicFields":[]},{"name":"text_cjk","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.CJKWidthFilterFactory"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.CJKBigramFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_ckb","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.SoraniNormalizationFilterFactory"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ckb.txt","ignoreCase":"true"},{"class":"solr.SoraniStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_cz","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_cz.txt","ignoreCase":"true"},{"class":"solr.CzechStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_da","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_da.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Danish"}]},"fields":[],"dynamicFields":[]},{"name":"text_de","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_de.txt","ignoreCase":"true"},{"class":"solr.GermanNormalizationFilterFactory"},{"class":"solr.GermanLightStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_el","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.GreekLowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_el.txt","ignoreCase":"false"},{"class":"solr.GreekStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_en","class":"solr.TextField","positionIncrementGap":"100","indexAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.EnglishPossessiveFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.PorterStemFilterFactory"}]},"queryAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.SynonymFilterFactory","expand":"true","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.EnglishPossessiveFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.PorterStemFilterFactory"}]},"fields":[],"dynamicFields":["*_en"]},{"name":"text_en_splitting","class":"solr.TextField","autoGeneratePhraseQueries":"true","positionIncrementGap":"100","indexAnalyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.WordDelimiterFilterFactory","catenateNumbers":"1","generateNumberParts":"1","splitOnCaseChange":"1","generateWordParts":"1","catenateAll":"0","catenateWords":"1"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.PorterStemFilterFactory"}]},"queryAnalyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"},"filters":[{"class":"solr.SynonymFilterFactory","expand":"true","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.WordDelimiterFilterFactory","catenateNumbers":"0","generateNumberParts":"1","splitOnCaseChange":"1","generateWordParts":"1","catenateAll":"0","catenateWords":"0"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.PorterStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_en_splitting_tight","class":"solr.TextField","autoGeneratePhraseQueries":"true","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"},"filters":[{"class":"solr.SynonymFilterFactory","expand":"false","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_en.txt","ignoreCase":"true"},{"class":"solr.WordDelimiterFilterFactory","catenateNumbers":"1","generateNumberParts":"0","generateWordParts":"0","catenateAll":"0","catenateWords":"1"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.KeywordMarkerFilterFactory","protected":"protwords.txt"},{"class":"solr.EnglishMinimalStemFilterFactory"},{"class":"solr.RemoveDuplicatesTokenFilterFactory"}]},"fields":["sku"],"dynamicFields":[]},{"name":"text_es","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_es.txt","ignoreCase":"true"},{"class":"solr.SpanishLightStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_eu","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_eu.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Basque"}]},"fields":[],"dynamicFields":[]},{"name":"text_fa","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"charFilters":[{"class":"solr.PersianCharFilterFactory"}],"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.ArabicNormalizationFilterFactory"},{"class":"solr.PersianNormalizationFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_fa.txt","ignoreCase":"true"}]},"fields":[],"dynamicFields":[]},{"name":"text_fi","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_fi.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Finnish"}]},"fields":[],"dynamicFields":[]},{"name":"text_fr","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.ElisionFilterFactory","articles":"lang/contractions_fr.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_fr.txt","ignoreCase":"true"},{"class":"solr.FrenchLightStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_ga","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.ElisionFilterFactory","articles":"lang/contractions_ga.txt","ignoreCase":"true"},{"class":"solr.StopFilterFactory","words":"lang/hyphenations_ga.txt","ignoreCase":"true"},{"class":"solr.IrishLowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ga.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Irish"}]},"fields":[],"dynamicFields":[]},{"name":"text_general","class":"solr.TextField","positionIncrementGap":"100","indexAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"stopwords.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"}]},"queryAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"stopwords.txt","ignoreCase":"true"},{"class":"solr.SynonymFilterFactory","expand":"true","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.LowerCaseFilterFactory"}]},"fields":["author","category","comments","content","description","features","includes","keywords","manu","name","resourcename","subject","text","title","url"],"dynamicFields":["attr_*","*_txt","*_t"]},{"name":"text_general_rev","class":"solr.TextField","positionIncrementGap":"100","indexAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.StopFilterFactory","words":"stopwords.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.ReversedWildcardFilterFactory","maxPosQuestion":"2","maxFractionAsterisk":"0.33","maxPosAsterisk":"3","withOriginal":"true"}]},"queryAnalyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.SynonymFilterFactory","expand":"true","ignoreCase":"true","synonyms":"synonyms.txt"},{"class":"solr.StopFilterFactory","words":"stopwords.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"}]},"fields":["text_rev"],"dynamicFields":[]},{"name":"text_gl","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_gl.txt","ignoreCase":"true"},{"class":"solr.GalicianStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_hi","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.IndicNormalizationFilterFactory"},{"class":"solr.HindiNormalizationFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_hi.txt","ignoreCase":"true"},{"class":"solr.HindiStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_hu","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_hu.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Hungarian"}]},"fields":[],"dynamicFields":[]},{"name":"text_hy","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_hy.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Armenian"}]},"fields":[],"dynamicFields":[]},{"name":"text_id","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_id.txt","ignoreCase":"true"},{"class":"solr.IndonesianStemFilterFactory","stemDerivational":"true"}]},"fields":[],"dynamicFields":[]},{"name":"text_it","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.ElisionFilterFactory","articles":"lang/contractions_it.txt","ignoreCase":"true"},{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_it.txt","ignoreCase":"true"},{"class":"solr.ItalianLightStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_ja","class":"solr.TextField","autoGeneratePhraseQueries":"false","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.JapaneseTokenizerFactory","mode":"search"},"filters":[{"class":"solr.JapaneseBaseFormFilterFactory"},{"class":"solr.JapanesePartOfSpeechStopFilterFactory","tags":"lang/stoptags_ja.txt"},{"class":"solr.CJKWidthFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ja.txt","ignoreCase":"true"},{"class":"solr.JapaneseKatakanaStemFilterFactory","minimumLength":"4"},{"class":"solr.LowerCaseFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_lv","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_lv.txt","ignoreCase":"true"},{"class":"solr.LatvianStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_nl","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_nl.txt","ignoreCase":"true"},{"class":"solr.StemmerOverrideFilterFactory","dictionary":"lang/stemdict_nl.txt","ignoreCase":"false"},{"class":"solr.SnowballPorterFilterFactory","language":"Dutch"}]},"fields":[],"dynamicFields":[]},{"name":"text_no","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_no.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Norwegian"}]},"fields":[],"dynamicFields":[]},{"name":"text_pt","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_pt.txt","ignoreCase":"true"},{"class":"solr.PortugueseLightStemFilterFactory"}]},"fields":[],"dynamicFields":[]},{"name":"text_ro","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_ro.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Romanian"}]},"fields":[],"dynamicFields":[]},{"name":"text_ru","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_ru.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Russian"}]},"fields":[],"dynamicFields":[]},{"name":"text_sv","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","format":"snowball","words":"lang/stopwords_sv.txt","ignoreCase":"true"},{"class":"solr.SnowballPorterFilterFactory","language":"Swedish"}]},"fields":[],"dynamicFields":[]},{"name":"text_th","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.LowerCaseFilterFactory"},{"class":"solr.ThaiWordFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_th.txt","ignoreCase":"true"}]},"fields":[],"dynamicFields":[]},{"name":"text_tr","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"},"filters":[{"class":"solr.TurkishLowerCaseFilterFactory"},{"class":"solr.StopFilterFactory","words":"lang/stopwords_tr.txt","ignoreCase":"false"},{"class":"solr.SnowballPorterFilterFactory","language":"Turkish"}]},"fields":[],"dynamicFields":[]},{"name":"text_ws","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.WhitespaceTokenizerFactory"}},"fields":[],"dynamicFields":[]},{"name":"tfloat","class":"solr.TrieFloatField","positionIncrementGap":"0","precisionStep":"8","fields":[],"dynamicFields":["*_tf"]},{"name":"tint","class":"solr.TrieIntField","positionIncrementGap":"0","precisionStep":"8","fields":[],"dynamicFields":["*_ti"]},{"name":"tlong","class":"solr.TrieLongField","positionIncrementGap":"0","precisionStep":"8","fields":[],"dynamicFields":["*_tl"]}]}`)
}

func mockSchemaFieldTypesLocation(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":2},"fieldType":{"name":"location","class":"solr.LatLonType","subFieldSuffix":"_coordinate","fields":["store"],"dynamicFields":["*_p"]}}`)
}

func mockSchemaName(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":1},"name":"example"}`)
}

func mockSchemaUniquekey(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":3},"uniqueKey":"id"}`)
}

func mockSchemaVersion(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":1},"version":1.5}`)
}

func writeContentTypeError(w http.ResponseWriter) {
	io.WriteString(w, `{
		  "error":{
		    "msg":"Must specify a Content-Type header with POST requests",
		    "code":415}}`)

}

func mockPing(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":2,"params":{"df":"text","echoParams":"all","rows":"10","echoParams":"all","wt":"json","ts":"1408264558581","_":"1408264558582","q":"solrpingquery","distrib":"false"}},"status":"OK"}`)
}

func mockMoreLikeThisSuccess(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":0,
		    "QTime":12},
		  "match":{"numFound":200,"start":0,"docs":[
		      {
		        "id":"test_id_0",
		        "title":["add sucess 0"],
		        "_version_":1476345720316362752}]
		  },
		  "response":{"numFound":199,"start":0,"docs":[
		      {
		        "id":"test_id_1",
		        "title":["add sucess 1"],
		        "_version_":1476345720644567040},
		      {
		        "id":"test_id_2",
		        "title":["add sucess 2"],
		        "_version_":1476345720645615616},
		      {
		        "id":"test_id_3",
		        "title":["add sucess 3"],
		        "_version_":1476345720645615617}]
		  }}`)
}

func mockMoreLikeThisError(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	w.WriteHeader(400)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":400,
		    "QTime":5},
		  "error":{
		    "msg":"Missing required parameter: mlt.fl",
		    "code":400}}`)
}




func mockStartServer() {
	http.HandleFunc("/success/core0/select/", mockSuccessSelect)
	http.HandleFunc("/fail/core0/select/", mockFailSelect)
	http.HandleFunc("/facet_counts/core0/select/", mockSuccessSelectFacet)
	http.HandleFunc("/highlight/core0/select/", mockSuccessSelectHighlight)

	http.HandleFunc("/standalonecommit/core0/update/", mockSuccessStandaloneCommit)
	http.HandleFunc("/add/core0/update/", mockSuccessAdd)
	http.HandleFunc("/delete/core0/update/", mockSuccessDelete)

	http.HandleFunc("/command/core0/update/", mockSuccessCommand)
	http.HandleFunc("/xml/core0/update/", mockSuccessXML)
	http.HandleFunc("/grouped/core0/select/", mockSuccessGrouped)
	http.HandleFunc("/noresponse/core0/select/", mockSuccessStrangeGrouped)
	http.HandleFunc("/solr/admin/cores", mockCoreAdmin)
	http.HandleFunc("/stats/collection1/select", mockSuccessStats)
	http.HandleFunc("/success/collection1/mlt", mockMoreLikeThisSuccess)
	http.HandleFunc("/error/collection1/mlt", mockMoreLikeThisError)
	
	http.HandleFunc("/solr/collection1/schema", mockSchema)
	http.HandleFunc("/solr/collection1/schema/fields", mockSchemaFields)
	http.HandleFunc("/solr/collection1/schema/fields/title", mockSchemaFieldsTitle)
	http.HandleFunc("/solr/collection1/schema/dynamicfields", mockSchemaDynamicFields)
	http.HandleFunc("/solr/collection1/schema/dynamicfields/*_coordinate", mockSchemaDynamicFieldsCoordinate)
	http.HandleFunc("/solr/collection1/schema/fieldtypes", mockSchemaFieldTypes)
	http.HandleFunc("/solr/collection1/schema/fieldtypes/location", mockSchemaFieldTypesLocation)
	http.HandleFunc("/solr/collection1/schema/name", mockSchemaName)
	http.HandleFunc("/solr/collection1/schema/uniquekey", mockSchemaUniquekey)
	http.HandleFunc("/solr/collection1/schema/version", mockSchemaVersion)
	http.HandleFunc("/solr/collection1/admin/ping", mockPing)
	http.HandleFunc("/solr/collection1/spell", mockSuccessSpell)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
