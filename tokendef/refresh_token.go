package tokendef

type OAuth2RefreshToken struct {
	*OAuth2StandardClaims
	SequenceNumber int `json:"seq"`
}
