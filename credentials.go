package qapi

const loginServerURL = "https://login.questrade.com/oauth2/"
const practiceLoginServerURL = "https://practicelogin.questrade.com/oauth2/"

type LoginCredentials struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ApiServer    string `json:"api_server"`
}

// authHeader is a shortcut that returns a string to be placed
// in the authorization header of API calls
func (l LoginCredentials) authHeader() string {
	return l.TokenType + " " + l.AccessToken
}
