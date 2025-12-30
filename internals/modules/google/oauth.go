package google

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

//store the access and referesh token in the memory

//check if referesh and access token already present
//skip getting oauth url if present

func GetOAuthUrl() string {
	//get the client_id from the env
	client_id := os.Getenv("OAUTH_CLIENT_ID")

	//set the base url to the google oauth endpoint
	u, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	if err != nil {
		log.Fatal("couldn't parse the base url")
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
	return u.String()
}

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

// get the referesh token and the access token
// also serve all the possible errors
func GetAccRefTokens(code string) (string, error) {
	//form the req url
	u, err := url.Parse("https://oauth2.googleapis.com/token")
	if err != nil {
		log.Fatal("cannot parse the base url : ", err.Error())
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

	//storing this in a json file, need to shift to a redis, and a way to tag the token with a particular user -todo later-
	var body map[string]any

	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		log.Fatal("error decoding the body : ", err.Error())
	}

	token_raw, err := os.ReadFile("./internals/modules/google/token.json")
	if err != nil {
		log.Fatal("error opening the token.json file :", err.Error())
	}

	var token Token
	err = json.Unmarshal(token_raw, &token)
	if err != nil {
		log.Fatal("error while unmarshaling the raw token data : ", err.Error())
	}

	fmt.Println(body["expires_in"]) //redundant fmt -delete later-

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

	token_raw, err = json.Marshal(token)
	if err != nil {
		log.Fatal("error marshaling the token : ", err.Error())
	}

	err = os.WriteFile("./internals/modules/google/token.json", token_raw, 0666)
	if err != nil {
		log.Fatal("error while writing the file :", err.Error())
	}
	return "success", nil

}

//write a function to referesh the token
