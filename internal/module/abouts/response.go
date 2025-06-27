package abouts

type AboutResponse struct {
	Name   string `json:"name"`
	Env    string `json:"env"`
	Locale string `json:"locale"`
}
