package inputForms

var (
	defaultPage    = 1
	defaultPerPage = 5
)

type GetData struct {
	Symbol  string `form:"symbol"`
	Open    string `form:"open"`
	Unix    string `form:"unix"`
	High    string `form:"high"`
	Low     string `form:"low"`
	Close   string `form:"close"`
	Page    int    `form:"page"`
	PerPage int    `form:"per_page"`
}

func (i *GetData) Handle() {
	if i.Page <= 0 {
		i.Page = defaultPage
	}

	if i.PerPage <= 0 {
		i.PerPage = defaultPerPage
	}
}
