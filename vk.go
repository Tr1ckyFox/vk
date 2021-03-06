package vk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

//const clientID = "2274003"                      //VK for Android app client_id
//const clientSecret = "hHbZxrka2uZ6jB1inYsH"     //VK for Android app client_secret
//const clientID = "3697615"                      //VK for Windows app client_id
//const clientSecret = "AlVXZFMUqyrnABp8ncuU"     //VK for Windows app client_secret
const clientID = "3140623"                        //VK for iPhone app client_id
const clientSecret = "VeWdmVclDCtn6ihuP1nt"       //VK for iPhone app client_secret
const authURL = "https://oauth.vk.com/token?"     //Direct Authorization URL
const apiMethodURL = "https://api.vk.com/method/" //Method request URL

// AuthResponse structure contains all parameters of response of authorization request
type AuthResponse struct {
	UserID           int    `json:"user_id"`
	ExpiresIn        int    `json:"expires_in"`
	AccessToken      string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Auth function makes authorization request and returns *AuthResponse structure
func Auth(login string, password string) (*AuthResponse, error) {
	var jsonResponse *AuthResponse
	var requestURL = url.Values{
		"grant_type":    {"password"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {login},
		"password":      {password},
	}
	response, err := http.Get(authURL + requestURL.Encode())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(content, &jsonResponse); err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

// Request function makes api method request and returns []byte JSON response
func Request(methodName string, parameters map[string]string, user *AuthResponse) ([]byte, error) {
	requestURL, err := url.Parse(apiMethodURL + methodName)
	if err != nil {
		return nil, err
	}
	requestQuery := requestURL.Query()
	for key, value := range parameters {
		requestQuery.Set(key, value)
	}
	if user != nil {
		requestQuery.Set("access_token", user.AccessToken)
	}
	requestURL.RawQuery = requestQuery.Encode()
	response, err := http.Get(requestURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}
