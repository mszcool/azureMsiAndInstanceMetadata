package main

import (
	"fmt"
	"log"
	"time"
 	"net/http"
)

/*Index (w, r)
 *Returns with a list of available functions for this simple API*/
 func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!");
}

/*MyMeta (w, r)
 *Returns instance metadata retrieved through the in-VM instance metadata service of the VM*/
func MyMeta(w http.ResponseWriter, r *http.Request) {
	metaDataJSON := GetInstanceMetadata()
	fmt.Fprintf(w, metaDataJSON)
}

/*MyPeers (w, r)
 *Uses the MSI to get a token and list all the other servers available in the resource group*/
func MyPeers(w http.ResponseWriter, r *http.Request) {
	token, err := GetMsiToken(50342)
	if err != "" {
		t := time.Now()
		log.Printf("--- %s --- Failed requesting MSI token --- %s", t.Format(time.RFC3339Nano), err)
		fmt.Fprint(w, err)
	} else {
		peerVms, err := GetMyPeerVirtualMachines(token)
		if err != "" {
			t := time.Now()
			log.Printf("--- %s --- Failed retrieving VMs from Azure Resource Manager APIs --- %s", t.Format(time.RFC3339Nano), err)
			fmt.Fprint(w, err)
		} else {
			fmt.Fprint(w, peerVms)
		}
	}
}