//	PhalGo-Json
//	Json解析功能,重写widuu/gojson新增部分功能
//	喵了个咪 <wenzhenxi@vip.qq.com> 2016/5/11
//  依赖情况: 无依赖
//  附上gojson地址:github.com/widuu/gojson

package phalgo

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Js struct {
	Data interface{}
}


//Initialize the json configruation
func Json(data string) *Js {
	j := new(Js)
	var f interface{}
	err := json.Unmarshal([]byte(data), &f)
	if err != nil {
		return j
	}
	j.Data = f
	return j
}



//According to the key of the returned data information,return js.data
func (j *Js) Get(key string) *Js {
	m := j.Getdata()
	if v, ok := m[key]; ok {
		j.Data = v
		return j
	}
	j.Data = nil
	return j
}

//return json data
func (j *Js) Getdata() map[string]interface{} {
	if m, ok := (j.Data).(map[string]interface{}); ok {
		return m
	}
	return nil
}

func (j *Js) Getindex(i int) *Js {

	num := i - 1
	if m, ok := (j.Data).([]interface{}); ok {
		v := m[num]
		j.Data = v
		return j
	}

	if m, ok := (j.Data).(map[string]interface{}); ok {
		var n = 0
		var data = make(map[string]interface{})
		for i, v := range m {
			if n == num {
				switch vv := v.(type) {
				case float64:
					data[i] = strconv.FormatFloat(vv, 'f', -1, 64)
					j.Data = data
					return j
				case string:
					data[i] = vv
					j.Data = data
					return j
				case []interface{}:
					j.Data = vv
					return j
				}

			}
			n++
		}

	}
	j.Data = nil
	return j
}

// When the data {"result":["username","password"]} can use arrayindex(1) get the username
func (j *Js) Arrayindex(i int) string {
	num := i - 1
	if i > len((j.Data).([]interface{})) {
		data := errors.New("index out of range list").Error()
		return data
	}
	if m, ok := (j.Data).([]interface{}); ok {
		v := m[num]
		switch vv := v.(type) {
		case float64:
			return strconv.FormatFloat(vv, 'f', -1, 64)
		case string:
			return vv
		default:
			return ""
		}

	}

	if _, ok := (j.Data).(map[string]interface{}); ok {
		return "error"
	}
	return "error"

}

//The data must be []interface{} ,According to your custom number to return your key and array data
func (j *Js) Getkey(key string, i int) *Js {
	num := i - 1
	if i > len((j.Data).([]interface{})) {
		j.Data = errors.New("index out of range list").Error()
		return j
	}
	if m, ok := (j.Data).([]interface{}); ok {
		v := m[num].(map[string]interface{})
		if h, ok := v[key]; ok {
			j.Data = h
			return j
		}

	}
	j.Data = nil
	return j
}

//According to the custom of the PATH to find the PATH
func (j *Js) Getpath(args ...string) *Js {
	d := j
	for i := range args {
		m := d.Getdata()

		if val, ok := m[args[i]]; ok {
			d.Data = val
		} else {
			d.Data = nil
			return d
		}
	}
	return d
}

//----------------------------------------新增-------------------------------------

func (j *Js)ToData() interface{} {
	return j.Data
}

func (j *Js)ToSlice() []interface{} {
	if m, ok := (j.Data).([]interface{}); ok {
		return m
	}
	return nil
}

func (j *Js) ToInt() int {
	if m, ok := j.Data.(int); ok {
		return m
	}
	return 0
}

func (j *Js) ToFloat() float64 {

	if m, ok := j.Data.(float64); ok {
		return m
	}
	if m, ok := j.Data.(string); ok {
		if m, ok := strconv.ParseFloat(m, 64); ok != nil {
			return m;
		}
	}
	return 0
}

//----------------------------------------新增----------------------------------------

func (j *Js) Tostring() string {
	if m, ok := j.Data.(string); ok {
		return m
	}
	if m, ok := j.Data.(float64); ok {
		return strconv.FormatFloat(m, 'f', -1, 64)
	}
	return ""
}

func (j *Js) ToArray() (k, d []string) {
	var key, data []string
	if m, ok := (j.Data).([]interface{}); ok {
		for _, value := range m {
			for index, v := range value.(map[string]interface{}) {
				switch vv := v.(type) {
				case float64:
					data = append(data, strconv.FormatFloat(vv, 'f', -1, 64))
					key = append(key, index)
				case string:
					data = append(data, vv)
					key = append(key, index)

				}
			}
		}

		return key, data
	}

	if m, ok := (j.Data).(map[string]interface{}); ok {
		for index, v := range m {
			switch vv := v.(type) {
			case float64:
				data = append(data, strconv.FormatFloat(vv, 'f', -1, 64))
				key = append(key, index)
			case string:
				data = append(data, vv)
				key = append(key, index)
			}
		}
		return key, data
	}

	return nil, nil
}

//获取[]string类型,整数和浮点类型会被转换成string
func (j *Js) StringtoArray() []string {
	var data []string

	switch j.Data.(type) {
	case []interface{}:
		for _, v := range j.Data.([]interface{}) {
			switch vv := v.(type) {
			case string:
				data = append(data, vv)
			case float64:
				data = append(data, strconv.FormatFloat(vv, 'f', -1, 64))
			}
		}
	}

	return data
}