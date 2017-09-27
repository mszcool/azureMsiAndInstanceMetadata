package main

import (
	"log"
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
)

const instanceMetaDataURL string = "http://169.254.169.254/metadata/instance?api-version=2017-04-02"

/*GetInstanceMetadata ()
 *Calls the Azure in-VM Instance Metadata service and returns the results to the caller*/
func GetInstanceMetadata() string {

	// Log the entry point
	t := time.Now()
	log.Printf("--- %s --- GetInstanceMetadata()", t.Format(time.RFC3339Nano))

	// Build a request to call the instance Azure in-VM metadata service
	req, err := http.NewRequest("GET", instanceMetaDataURL, nil)
	if err != nil {
		log.Printf("--- %s --- Failed creating http request --- %s", t.Format(time.RFC3339Nano), err)
		return "{ \"error\": \"failed creating http request object to retrieve instance metadata!\" }"
	}

	// Set the required header for the HTTP request
	req.Header.Add("Metadata", "true")

	// Create the HTTP client and call the instance metadata service
	client := &http.Client{}
	resp, err := client.Do(req);
	if err != nil {
		t = time.Now()
		log.Printf("--- %s --- Failed calling instance metadata service --- %s", t.Format(time.RFC3339Nano), err)
		return "{ \"error\": \"failed calling the in-VM instance metadata service!\" }"
	}
	// Complete reading the body
	defer resp.Body.Close()

	// Now return the instance metadata JSON or another error if the status code is not in 2xx range
	if (resp.StatusCode >= 200) && (resp.StatusCode <= 299) {
		bodyContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t = time.Now()
			log.Printf("--- %s --- Failed reading results from instance metadata service --- %s", t.Format(time.RFC3339Nano), err)
			return "{ \"error\": \"failed reading results from instance metadata service!\" }"
		}

		t = time.Now()
		log.Printf("--- %s --- Succeeded", t.Format(time.RFC3339Nano))

		return string(bodyContent)
	}

	t = time.Now()
	log.Printf("--- %s --- Failed with Non-200 status code: %q", t.Format(time.RFC3339Nano), resp.StatusCode)

	return fmt.Sprintf("{ \"error\": \"instance meta data service returned non-OK status code: %q \" }", resp.StatusCode)
}