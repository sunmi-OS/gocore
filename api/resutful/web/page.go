package web

type Pager struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

func (p Pager) Apply(total int, data interface{}) interface{} {
	if data == nil {
		return nil
	}
	res := new(struct {
		List interface{} `json:"list"`
		Page struct {
			Total    int `json:"total"`
			PageNo   int `json:"page_no"`
			PageSize int `json:"page_size"`
		} `json:"page"`
	})
	res.List = data
	res.Page.PageNo = p.PageNo
	res.Page.PageSize = p.PageSize
	res.Page.Total = total
	return res
}
