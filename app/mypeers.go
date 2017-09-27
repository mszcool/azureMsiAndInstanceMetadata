package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"net/http"
	"io/ioutil"
)

const (
	// Resource Group and Subscription ID should be set through environment variables in the system
	environmentNameSubscription string = "SUBSCRIPTION_ID"
	environmentNameResourceGroup string = "RESOURCE_GROUP"

	// Endpoints for the Azure REST API
	// Note: at the time of writing, Managed Service Identities were not fully supported in the GO-SDK, yet. 
	//       Also, I found this approach demonstrates the concepts of MSIs even more explicit compared to just relying on the SDK.
	restAPIEndpoint string = "https://management.azure.com/subscriptions/%s/resourceGroups/%s/%s"
	vmRelativeEndpoint string = "providers/Microsoft.Compute/virtualmachines?api-version=2016-04-30-preview"
	authorizationHeader string = "%s %s"
)

/*GetMyPeerVirtualMachines (msiToken)
 *Returns all the virtual machines in the same resource group by using the Azure Resource Manager
 *REST APIs. For authenticating against the REST APIs, it's using the token passed into the function*/
func GetMyPeerVirtualMachines(msiToken MsiToken) (vms string, errOut string) {

	// Log the entry point
	t := time.Now()
	log.Printf("--- %s --- GetMyPeerVirtualMachines()", t.Format(time.RFC3339Nano))

	// Then retrieve all required environment variables
	subID := os.Getenv(environmentNameSubscription)
	resGroup := os.Getenv(environmentNameResourceGroup)
	if (subID == "") || (resGroup == "") {
		log.Printf("--- %s --- Either the environment variable %s or %s is not set on the system!", t.Format(time.RFC3339Nano), environmentNameSubscription, environmentNameResourceGroup)
		return "", "Missing Environment variables for reading Virtual Machine peers in Resource Group!"
	}

	// Create the final endpoint URLs to call into the Azure Resource Manager VM REST API
	finalURL := fmt.Sprintf(restAPIEndpoint, subID, resGroup, vmRelativeEndpoint)
	finalAuthHeader := fmt.Sprintf(authorizationHeader, msiToken.TokenType, msiToken.AccessToken)

	// Build a request to call the instance Azure in-VM metadata service
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		log.Printf("--- %s --- Failed creating http request --- %s", t.Format(time.RFC3339Nano), err)
		return "", "{ \"error\": \"failed creating http request object to call Azure RM APIs!\" }"
	}
	req.Header.Add("Authorization", finalAuthHeader)

	// Create the HTTP client and call the instance metadata service
	client := &http.Client{}
	resp, err := client.Do(req);
	if err != nil {
		t = time.Now()
		log.Printf("--- %s --- Failed calling Azure Resource Manager REST API --- %s", t.Format(time.RFC3339Nano), err)
		return "", "{ \"error\": \"failed calling the Azure Resource Manager REST API!\" }"
	}
	// Complete reading the body
	defer resp.Body.Close()

	// Now return the raw VM JSON or another error if the status code is not in 2xx range
	if (resp.StatusCode >= 200) && (resp.StatusCode <= 299) {
		bodyContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t = time.Now()
			log.Printf("--- %s --- Failed reading VM JSON from Azure Resource Manager REST API --- %s", t.Format(time.RFC3339Nano), err)
			return "", "{ \"error\": \"failed reading VM JSON from Azure Resource Manager REST API!\" }"
		}

		t = time.Now()
		log.Printf("--- %s --- Succeeded", t.Format(time.RFC3339Nano))

		return string(bodyContent), ""
	}

	t = time.Now()
	log.Printf("--- %s --- Failed with Non-200 status code: %d", t.Format(time.RFC3339Nano), resp.StatusCode)

	return "", fmt.Sprintf("{ \"error\": \"Azure Resource Manager REST API call returned non-OK status code: %d \" }", resp.StatusCode)
}