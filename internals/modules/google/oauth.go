package google

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"
)

// ----------------------
// store the access and referesh token in the memory
// load func
func LoadAccRefToken(token *Token) error {
	token_raw, err := os.ReadFile("./internals/modules/google/token.json")
	if err != nil {
		return errors.Join(errors.New("error opening the token.json file"), err)
	}

	err = json.Unmarshal(token_raw, token)
	if err != nil {
		return errors.Join(errors.New("error while unmarshaling the raw token data : "), err)
	}

	return nil
}

//----------------------

// ----------------------
// save func
func SaveAccRefToken(token *Token) error {
	token_raw, err := json.Marshal(token)
	if err != nil {
		return errors.Join(errors.New("error while marshaling the token"), err)
	}

	err = os.WriteFile("./internals/modules/google/token.json", token_raw, 0666)
	if err != nil {
		return errors.Join(errors.New("error while writing the file"), err)
	}

	return nil
}

//----------------------

// ----------------------
// func to check if referesh token has expired
func IsRefTokenValid(token *Token) bool {
	//check if expire time is greater than time now
	if token.RefreshTokenExpiresIn-60 > time.Now().Unix() { // a 60 sec buffer to avoid error
		return true
	}
	return false
}

//----------------------

// ----------------------
// check if referesh and access token already present
// skip getting oauth url if present
func GetOAuthUrl() (string, error) {
	//get the client_id from the env
	client_id := os.Getenv("OAUTH_CLIENT_ID")

	//set the base url to the google oauth endpoint
	u, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	if err != nil {
		return "", errors.Join(errors.New("couldn't parse the base url"), err)
	}

	//set the required and recommeded queries on the base url
	q := u.Query()
	q.Add("client_id", client_id)
	q.Add("redirect_uri", "http://localhost:3030/google/code_redirect") //change this to the endpoint which serves the redirect req, and add the redirect uri on google cloud
	q.Add("response_type", "code")
	q.Add("scope", "https://www.googleapis.com/auth/youtube.upload")
	q.Add("access_type", "offline")

	//form the final url string
	u.RawQuery = q.Encode()

	//return the string
	return u.String(), nil
}

//----------------------

// ----------------------
// get the auth code from the redirected url
// also serve all the possible errors
func GetAuthCode(u *url.URL) (string, error) {
	//get the code from the queries
	q := u.Query()

	//serve if error occured
	if q.Get("error") == "access_denied" {
		return "", errors.New("access_denied")
	}

	c := q.Get("code")
	return c, nil
}

//----------------------

// ----------------------
// get the referesh token and the access token
// also serve all the possible errors
func GetAccRefTokens(code string, token *Token) error {
	//form the req url
	u, err := url.Parse("https://oauth2.googleapis.com/token")
	if err != nil {
		return errors.Join(errors.New("cannot parse the base url"), err)
	}

	q := u.Query()
	q.Add("client_id", os.Getenv("OAUTH_CLIENT_ID"))
	q.Add("client_secret", os.Getenv("OAUTH_CLIENT_SECRET"))
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	q.Add("redirect_uri", "http://localhost:3030/google/code_redirect")

	u.RawQuery = q.Encode()

	//make the post req
	res, err := http.Post(u.String(), "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		return errors.Join(errors.New("error while making the post req"), err)
	}

	//storing this in a json file, need to shift to a redis, and a way to tag the token with a particular user -todo later-
	var body map[string]any

	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return errors.Join(errors.New("error while decoding the response body"), err)
	}

	//update the token values
	if _, ok := body["access_token"].(string); ok {
		token.AccessToken = body["access_token"].(string)
	}
	if _, ok := body["refresh_token"].(string); ok {
		token.RefreshToken = body["refresh_token"].(string)
	}
	if _, ok := body["expires_in"].(float64); ok {
		token.ExpiresIn = time.Now().Unix() + int64(body["expires_in"].(float64))
	}
	if _, ok := body["refresh_token_expires_in"].(float64); ok {
		token.RefreshTokenExpiresIn = time.Now().Unix() + int64(body["refresh_token_expires_in"].(float64))
	}

	err = SaveAccRefToken(token)
	if err != nil {
		return errors.Join(errors.New("error while saving the token"), err)
	}

	return nil // no error occured
}

//----------------------

// ----------------------
// func to check if acess token is valid
func IsAccTokenValid(token *Token) bool {
	// check if acc token expire time is still valid
	if token.ExpiresIn-60 > time.Now().Unix() { // a 60 sec buffer to avoid error
		return true
	}
	return false
}

//----------------------

// ----------------------
// write a function to referesh the token
func RefreshAccToken(token *Token) error {
	u, err := url.Parse("https://oauth2.googleapis.com/token") //setting the base url
	if err != nil {
		return errors.Join(errors.New("error setting the base url"), err)
	}

	q := u.Query()
	q.Add("client_id", os.Getenv("OAUTH_CLIENT_ID"))
	q.Add("client_secret", os.Getenv("OAUTH_CLIENT_SECRET"))
	q.Add("refresh_token", token.RefreshToken)
	q.Add("grant_type", "refresh_token")

	u.RawQuery = q.Encode()

	res, err := http.Post(u.String(), "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		return errors.Join(errors.New("error while making post req"), err)
	}

	var body map[string]any

	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return errors.Join(errors.New("error while decoding the response body"), err)
	}

	//update the token values
	if _, ok := body["access_token"].(string); ok {
		token.AccessToken = body["access_token"].(string)
	}
	if _, ok := body["expires_in"].(float64); ok {
		token.ExpiresIn = time.Now().Unix() + int64(body["expires_in"].(float64))
	}

	err = SaveAccRefToken(token)
	if err != nil {
		return errors.Join(errors.New("error while saving the token"), err)
	}

	return nil //no error occured
}

//----------------------

// ----------------------
// prepare function for the whole pipeline
func Prepare(token *Token) error {

	//load the token
	if *token == (Token{}) {
		err := LoadAccRefToken(token)
		if err != nil {
			return errors.Join(errors.New("error loading the access and refresh token"), err)
		}
	}

	//check if ref token valid
	if ok := IsRefTokenValid(token); !ok {
		return errors.New("ref token gone invalid, redo the auth procedure")
	}

	//check if acc token is valid
	if ok := IsAccTokenValid(token); !ok {
		err := RefreshAccToken(token)
		if err != nil {
			return errors.Join(errors.New("error refreshing the access token"), err)
		}
	}

	return nil
}

//----------------------
