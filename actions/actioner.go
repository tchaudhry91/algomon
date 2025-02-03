package actions

type Input struct {
	Params map[string]string `json:"params"`
}

type Actioner interface {
	Action(Input) error
}
