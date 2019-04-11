package solr

import "testing"

func BenchmarkStandardResultParser(b *testing.B) {
	data := []byte(`{"responseHeader":{"status":0,"QTime":0,"params":{"start":"0","q":"*:*","wt":"json","rows":"5"}},"response":{"numFound":12508577,"start":0,"docs":[{"id":"P27665908","p_name":"Archer Lady Case for Blackberry Q10","p_price_c":"175000.00,IDR","p_shop_name":"Nisa Toko","p_domain":"nisatoko","p_shop_sd_name":"Bekasi","p_key":"archer-lady-case-for-blackberry-q10","p_file_path":"product-1/2015/12/28/6208975","p_file_name":"6208975_3959f656-cb27-43df-83db-4b1c1455d3ce.jpg","p_id":27665908,"p_child_cat_id":69,"p_shop_id":733718,"p_condition":1,"p_server_id":1,"p_price":175000.0,"p_dep_id":["65","69","66"],"p_best_match_dink":-288.0,"p_best_match_price":4.756962,"_version_":1521775745340276736},{"id":"P27665917","p_name":"shampo zoku penumbuh rambut hitam berkilau","p_price_c":"115000.00,IDR","p_shop_name":"penumbuh rambut sehat10","p_domain":"penumbuhrambut93","p_shop_sd_name":"Surabaya","p_key":"shampo-zoku-penumbuh-rambut-hitam-berkilau","p_file_path":"product-1/2015/12/28/6766154","p_file_name":"6766154_3344bc2d-8069-4ef6-adc7-7de90b6604f5.jpg","p_id":27665917,"p_child_cat_id":436,"p_shop_id":801110,"p_condition":1,"p_server_id":1,"p_price":115000.0,"p_dep_id":["433","61","436"],"p_best_match_dink":-288.0,"p_best_match_price":4.939302,"_version_":1521775745341325313},{"id":"P27665953","p_name":"Archer Lady Case for HTC One M7","p_price_c":"175000.00,IDR","p_shop_name":"Nisa Toko","p_domain":"nisatoko","p_shop_sd_name":"Bekasi","p_key":"archer-lady-case-for-htc-one-m7","p_file_path":"product-1/2015/12/28/6208975","p_file_name":"6208975_68882b13-3e53-431b-bd64-d2e518db03e9.jpg","p_id":27665953,"p_child_cat_id":69,"p_shop_id":733718,"p_condition":1,"p_server_id":1,"p_price":175000.0,"p_dep_id":["65","69","66"],"p_best_match_dink":-288.0,"p_best_match_price":4.756962,"_version_":1521775745341325315},{"id":"P27665898","p_name":"Archer Lady Case for Blackberry 9900 9980","p_price_c":"175000.00,IDR","p_shop_name":"Nisa Toko","p_domain":"nisatoko","p_shop_sd_name":"Bekasi","p_key":"archer-lady-case-for-blackberry-9900-9980","p_file_path":"product-1/2015/12/28/6208975","p_file_name":"6208975_b40a60e6-1909-4a13-9d49-f881d3a8c5da.jpg","p_id":27665898,"p_child_cat_id":69,"p_shop_id":733718,"p_condition":1,"p_server_id":1,"p_price":175000.0,"p_dep_id":["65","69","66"],"p_best_match_dink":-288.0,"p_best_match_price":4.756962,"_version_":1521775745342373888},{"id":"P27665910","p_name":"VAMPIRE SERUM 30ml ORIGINAL","p_price_c":"26600.00,IDR","p_shop_name":"Dokter Cantikku","p_domain":"doktercantiku","p_shop_sd_name":"Jakarta","p_key":"vampire-serum-30ml-original","p_file_path":"product-1/2015/12/28/6634788","p_file_name":"6634788_e34791af-a170-48e8-b836-68da4c5d775a.jpg","p_id":27665910,"p_child_cat_id":598,"p_shop_id":785036,"p_condition":1,"p_server_id":1,"p_price":26600.0,"p_dep_id":["61","445","598"],"p_best_match_dink":-288.0,"p_best_match_price":5.5751185,"_version_":1521775745342373889}]}}`)

	for i := 0; i < b.N; i++ {

		parser := StandardResultParser{}

		parser.Parse(&data)

	}
}

func BenchmarkFireworkResultParser(b *testing.B) {
	data := []byte(`{"responseHeader":{"status":0,"QTime":0,"params":{"start":"0","q":"*:*","wt":"json","rows":"5"}},"response":{"numFound":12508577,"start":0,"docs":[{"id":"P27665908","p_name":"Archer Lady Case for Blackberry Q10","p_price_c":"175000.00,IDR","p_shop_name":"Nisa Toko","p_domain":"nisatoko","p_shop_sd_name":"Bekasi","p_key":"archer-lady-case-for-blackberry-q10","p_file_path":"product-1/2015/12/28/6208975","p_file_name":"6208975_3959f656-cb27-43df-83db-4b1c1455d3ce.jpg","p_id":27665908,"p_child_cat_id":69,"p_shop_id":733718,"p_condition":1,"p_server_id":1,"p_price":175000.0,"p_dep_id":["65","69","66"],"p_best_match_dink":-288.0,"p_best_match_price":4.756962,"_version_":1521775745340276736},{"id":"P27665917","p_name":"shampo zoku penumbuh rambut hitam berkilau","p_price_c":"115000.00,IDR","p_shop_name":"penumbuh rambut sehat10","p_domain":"penumbuhrambut93","p_shop_sd_name":"Surabaya","p_key":"shampo-zoku-penumbuh-rambut-hitam-berkilau","p_file_path":"product-1/2015/12/28/6766154","p_file_name":"6766154_3344bc2d-8069-4ef6-adc7-7de90b6604f5.jpg","p_id":27665917,"p_child_cat_id":436,"p_shop_id":801110,"p_condition":1,"p_server_id":1,"p_price":115000.0,"p_dep_id":["433","61","436"],"p_best_match_dink":-288.0,"p_best_match_price":4.939302,"_version_":1521775745341325313},{"id":"P27665953","p_name":"Archer Lady Case for HTC One M7","p_price_c":"175000.00,IDR","p_shop_name":"Nisa Toko","p_domain":"nisatoko","p_shop_sd_name":"Bekasi","p_key":"archer-lady-case-for-htc-one-m7","p_file_path":"product-1/2015/12/28/6208975","p_file_name":"6208975_68882b13-3e53-431b-bd64-d2e518db03e9.jpg","p_id":27665953,"p_child_cat_id":69,"p_shop_id":733718,"p_condition":1,"p_server_id":1,"p_price":175000.0,"p_dep_id":["65","69","66"],"p_best_match_dink":-288.0,"p_best_match_price":4.756962,"_version_":1521775745341325315},{"id":"P27665898","p_name":"Archer Lady Case for Blackberry 9900 9980","p_price_c":"175000.00,IDR","p_shop_name":"Nisa Toko","p_domain":"nisatoko","p_shop_sd_name":"Bekasi","p_key":"archer-lady-case-for-blackberry-9900-9980","p_file_path":"product-1/2015/12/28/6208975","p_file_name":"6208975_b40a60e6-1909-4a13-9d49-f881d3a8c5da.jpg","p_id":27665898,"p_child_cat_id":69,"p_shop_id":733718,"p_condition":1,"p_server_id":1,"p_price":175000.0,"p_dep_id":["65","69","66"],"p_best_match_dink":-288.0,"p_best_match_price":4.756962,"_version_":1521775745342373888},{"id":"P27665910","p_name":"VAMPIRE SERUM 30ml ORIGINAL","p_price_c":"26600.00,IDR","p_shop_name":"Dokter Cantikku","p_domain":"doktercantiku","p_shop_sd_name":"Jakarta","p_key":"vampire-serum-30ml-original","p_file_path":"product-1/2015/12/28/6634788","p_file_name":"6634788_e34791af-a170-48e8-b836-68da4c5d775a.jpg","p_id":27665910,"p_child_cat_id":598,"p_shop_id":785036,"p_condition":1,"p_server_id":1,"p_price":26600.0,"p_dep_id":["61","445","598"],"p_best_match_dink":-288.0,"p_best_match_price":5.5751185,"_version_":1521775745342373889}]}}`)

	for i := 0; i < b.N; i++ {
		parser := FireworkResultParser{}

		parser.Parse(&data)
	}
}

func BenchmarkExtensiveResultParser(b *testing.B) {
	data := []byte(`{"responseHeader":{"status":0,"QTime":0,"params":{"start":"0","q":"*:*","wt":"json","rows":"5"}},"response":{"numFound":12508577,"start":0,"docs":[{"id":"P27665908","p_name":"Archer Lady Case for Blackberry Q10","p_price_c":"175000.00,IDR","p_shop_name":"Nisa Toko","p_domain":"nisatoko","p_shop_sd_name":"Bekasi","p_key":"archer-lady-case-for-blackberry-q10","p_file_path":"product-1/2015/12/28/6208975","p_file_name":"6208975_3959f656-cb27-43df-83db-4b1c1455d3ce.jpg","p_id":27665908,"p_child_cat_id":69,"p_shop_id":733718,"p_condition":1,"p_server_id":1,"p_price":175000.0,"p_dep_id":["65","69","66"],"p_best_match_dink":-288.0,"p_best_match_price":4.756962,"_version_":1521775745340276736},{"id":"P27665917","p_name":"shampo zoku penumbuh rambut hitam berkilau","p_price_c":"115000.00,IDR","p_shop_name":"penumbuh rambut sehat10","p_domain":"penumbuhrambut93","p_shop_sd_name":"Surabaya","p_key":"shampo-zoku-penumbuh-rambut-hitam-berkilau","p_file_path":"product-1/2015/12/28/6766154","p_file_name":"6766154_3344bc2d-8069-4ef6-adc7-7de90b6604f5.jpg","p_id":27665917,"p_child_cat_id":436,"p_shop_id":801110,"p_condition":1,"p_server_id":1,"p_price":115000.0,"p_dep_id":["433","61","436"],"p_best_match_dink":-288.0,"p_best_match_price":4.939302,"_version_":1521775745341325313},{"id":"P27665953","p_name":"Archer Lady Case for HTC One M7","p_price_c":"175000.00,IDR","p_shop_name":"Nisa Toko","p_domain":"nisatoko","p_shop_sd_name":"Bekasi","p_key":"archer-lady-case-for-htc-one-m7","p_file_path":"product-1/2015/12/28/6208975","p_file_name":"6208975_68882b13-3e53-431b-bd64-d2e518db03e9.jpg","p_id":27665953,"p_child_cat_id":69,"p_shop_id":733718,"p_condition":1,"p_server_id":1,"p_price":175000.0,"p_dep_id":["65","69","66"],"p_best_match_dink":-288.0,"p_best_match_price":4.756962,"_version_":1521775745341325315},{"id":"P27665898","p_name":"Archer Lady Case for Blackberry 9900 9980","p_price_c":"175000.00,IDR","p_shop_name":"Nisa Toko","p_domain":"nisatoko","p_shop_sd_name":"Bekasi","p_key":"archer-lady-case-for-blackberry-9900-9980","p_file_path":"product-1/2015/12/28/6208975","p_file_name":"6208975_b40a60e6-1909-4a13-9d49-f881d3a8c5da.jpg","p_id":27665898,"p_child_cat_id":69,"p_shop_id":733718,"p_condition":1,"p_server_id":1,"p_price":175000.0,"p_dep_id":["65","69","66"],"p_best_match_dink":-288.0,"p_best_match_price":4.756962,"_version_":1521775745342373888},{"id":"P27665910","p_name":"VAMPIRE SERUM 30ml ORIGINAL","p_price_c":"26600.00,IDR","p_shop_name":"Dokter Cantikku","p_domain":"doktercantiku","p_shop_sd_name":"Jakarta","p_key":"vampire-serum-30ml-original","p_file_path":"product-1/2015/12/28/6634788","p_file_name":"6634788_e34791af-a170-48e8-b836-68da4c5d775a.jpg","p_id":27665910,"p_child_cat_id":598,"p_shop_id":785036,"p_condition":1,"p_server_id":1,"p_price":26600.0,"p_dep_id":["61","445","598"],"p_best_match_dink":-288.0,"p_best_match_price":5.5751185,"_version_":1521775745342373889}]}}`)

	for i := 0; i < b.N; i++ {
		parser := ExtensiveResultParser{}

		parser.Parse(&data)
	}
}

func TestParseMoreLikeThisMatch(t *testing.T) {
	data := []byte(`{
			  "responseHeader":{
			    "status":0,
			    "QTime":4},
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

	parser := MoreLikeThisParser{}

	res, _ := parser.Parse(&data)

	if res.Match.Start != 0 {
		t.Errorf("res.Match.Start expected to be '0' but got '%d'", res.Match.Start)
	}

	if len(res.Match.Docs) != 1 {
		t.Errorf("res.Match.Docs should have '1' doc but got '%d'", len(res.Match.Docs))
	}
	expected := "test_id_0"
	if res.Match.Docs[0].Get("id") != expected {
		t.Errorf("title expected to be '%s' but got '%s'", expected, res.Match.Docs[0].Get("id"))
	}

	if res.Results.Start != 0 {
		t.Errorf("res.Match.Start expected to be '0' but got '%d'", res.Results.Start)
	}

	if len(res.Results.Docs) != 3 {
		t.Errorf("res.Results.Docs should have '3' doc but got '%d'", len(res.Results.Docs))
	}
	expected = "test_id_1"
	if res.Results.Docs[0].Get("id") != expected {
		t.Errorf("title expected to be '%s' but got '%s'", expected, res.Results.Docs[0].Get("id"))
	}
}

func TestParseNextCursorMarkForFireworkResultParser(t *testing.T) {
	data := []byte(`{
	        "responseHeader":{
	          "status":0,
	          "QTime":4},
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
	        },
	        "nextCursorMark": "35665"
        }`)

	parser := FireworkResultParser{}
	res, _ := parser.Parse(&data)

	if res.NextCursorMark != "35665" {
		t.Errorf("Error")
	}
}

func TestParseNextCursorMarkForExtensiveResultParser(t *testing.T) {
	data := []byte(`{
	        "responseHeader":{
	          "status":0,
	          "QTime":4},
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
	        },
	        "nextCursorMark": "35665"
        }`)

	parser := ExtensiveResultParser{}
	res, _ := parser.Parse(&data)

	if res.NextCursorMark != "35665" {
		t.Errorf("Error")
	}
}

func TestParseNextCursorMarkForStandardResultParser(t *testing.T) {
	data := []byte(`{
	        "responseHeader":{
	          "status":0,
	          "QTime":4},
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
	        },
	        "nextCursorMark": "35665"
        }`)

	parser := StandardResultParser{}
	res, _ := parser.Parse(&data)

	if res.NextCursorMark != "35665" {
		t.Errorf("Error")
	}
}
