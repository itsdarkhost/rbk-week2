package country

type Country struct {
	Name   string   `json:"name"`
	Cities []string `json:"cities"`
}

type countriesResponse struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
	Data  []struct {
		Country string   `json:"country"`
		Cities  []string `json:"cities"`
	} `json:"data"`
}

type metadataResponse []struct {
	CCA2 string `json:"cca2"`
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
}
