package tokendef

type OAuth2AccessToken struct {
	*OAuth2StandardClaims
	DbCheck bool   `json:"dbchk"`
	Scope   string `json:"scope"`
	Role    string `json:"role"`
}
