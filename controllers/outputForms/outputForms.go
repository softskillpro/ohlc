package outputForms

type State struct {
	Message *string `json:"message"`
	Status  bool    `json:"status"`
	Code    int     `json:"code"`
	Detail  *string `json:"detail"`
	Counts  *Counts `json:"counts"`
	Data    any     `json:"data"`
}

type Counts struct {
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
}

func NewState() *State {
	s := new(State)

	return s
}

func (s *State) SetMessage(m string) *State {
	s.Message = &m
	return s
}

func (s *State) SetCode(c int) *State {
	s.Code = c
	return s
}

func (s *State) SetStatus(status bool) *State {
	s.Status = status
	return s
}

func (s *State) SetDetail(m string) *State {
	s.Detail = &m
	return s
}

func (s *State) SetCounts(m Counts) *State {
	s.Counts = &m
	return s
}

func (s *State) SetData(m any) *State {
	s.Data = m
	return s
}
