package main

import (
	"net/url"
	"io/ioutil"
	"log"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
)

const msiTokenURL string = "http://localhost:%d/oauth2/token"
const resourceURL string = "https://management.azure.com/"

/*MsiToken ()
 *Encapsulates a token retrieved by the MSI extension on the local machine*/
type MsiToken struct {
	AccessToken 	string `json:"access_token"`
	RefreshToken 	string `json:"refresh_token"`
	ExpiresIn 		string `json:"expires_in"`
	ExpiresOn 		string `json:"expires_on"`
	NotBefore 		string `json:"not_before"`
	Resource 		string `json:"resource"`
	TokenType 		string `json:"token_type"`
}

/*GetMsiToken (msiPort)
 *Uses the Managed Service Identity Extension to retrieve a token that allows the VM to call into
 *the Azure Resource Manager APIs*/
func GetMsiToken(msiPort int) (token MsiToken, errRet string) {

	var myToken MsiToken

	// Log the entry point
	t := time.Now()
	log.Printf("--- %s --- GetMsiToken()", t.Format(time.RFC3339Nano))

	// Build a request to call the MSI Extension OAuth2 Service
	// The request must contain the resource for which we request the token
	finalRequestURL := fmt.Sprintf("%s?resource=%s", fmt.Sprintf(msiTokenURL, msiPort), url.QueryEscape(resourceURL))
	req, err := http.NewRequest("GET", finalRequestURL, nil)
	if err != nil {
		log.Printf("--- %s --- Failed creating http request --- %s", t.Format(time.RFC3339Nano), err)
		return myToken, "{ \"error\": \"failed creating http request object to request MSI token!\" }"
	}

	// Set the required header for the HTTP request
	req.Header.Add("Metadata", "true")

	// Create the HTTP client and call the instance metadata service
	client := &http.Client{}
	resp, err := client.Do(req);
	if err != nil {
		t = time.Now()
		log.Printf("--- %s --- Failed calling MSI token service --- %s", t.Format(time.RFC3339Nano), err)
		return myToken, "{ \"error\": \"failed calling MSI token service!\" }"
	}
	// Complete reading the body
	defer resp.Body.Close()

	// Now return the instance metadata JSON or another error if the status code is not in 2xx range
	if (resp.StatusCode >= 200) && (resp.StatusCode <= 299) {
		dec := json.NewDecoder(resp.Body)
		err := dec.Decode(&myToken)
		if err != nil {
			t = time.Now()
			log.Printf("--- %s --- Failed decoding MSI token from MSI token endpoint --- %s", t.Format(time.RFC3339Nano), err)
			return myToken, "{ \"error\": \"failed decoding MSI token from MSI token endpoint!\" }"
		}

		t = time.Now()
		log.Printf("--- %s --- Succeeded", t.Format(time.RFC3339Nano))

		return myToken, ""
	}

	t = time.Now()
	log.Printf("--- %s --- Failed with Non-200 status code: %d", t.Format(time.RFC3339Nano), resp.StatusCode)

	// Try to read the body and log the error details, nevertheless
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("--- %s --- Failed reading response body from http response for more error details", t.Format(time.RFC3339Nano))
	} else {
		log.Printf("--- %s --- %s", t.Format(time.RFC3339Nano), string(bodyContent))
	}

	return myToken, fmt.Sprintf("{ \"error\": \"instance meta data service returned non-OK status code: %d \" }", resp.StatusCode)
}