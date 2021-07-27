package web

type Pager struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

type pageRsp struct {
	List interface{} `json:"list"`
	Page *page       `json:"page"`
}

type page struct {
	Total    int `json:"total"`
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

func (p Pager) Apply(total int, data interface{}) interface{} {
	if data == nil {
		return nil
	}
	res := &pageRsp{
		List: data,
		Page: &page{
			Total:    total,
			PageNo:   p.PageNo,
			PageSize: p.PageSize,
		},
	}
	return res
}
