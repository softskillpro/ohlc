package ohcl

const TableName = "ohcl"

type Model struct {
	Unix   int64   `sql:"unix" json:"unix"`
	Symbol string  `sql:"symbol" json:"symbol"`
	Open   float64 `sql:"open" json:"open"`
	High   float64 `sql:"high" json:"high"`
	Low    float64 `sql:"low" json:"low"`
	Close  float64 `sql:"close" json:"close"`
}

func (m Model) TableName() string {
	return TableName
}

func (r Model) GetRecords() []Model {
	return []Model{r}
}
