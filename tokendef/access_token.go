package tokendef

type OAuth2AccessToken struct {
	*OAuth2StandardClaims
	Scope string `json:"scope"`
}
