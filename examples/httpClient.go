// +build examples

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

const (
	port = ":9080"
)

func main() {
	l := log.New(os.Stderr, "GOKRB5 Client: ", log.LstdFlags)

	keytabBody, err := ioutil.ReadFile("./client.keytab")
	if err != nil {
		l.Fatalf("could not load client keytab: %v", err)
	}

	//defer profile.Start(profile.TraceProfile).Stop()
	// Load the keytab
	kt := keytab.New()
	err = kt.Unmarshal(keytabBody)
	if err != nil {
		l.Fatalf("could not load client keytab: %v", err)
	}

	kdc, err := ioutil.ReadFile("./krb5.conf")
	if err != nil {
		l.Fatalf("could not load krb5.conf: %v", err)
	}

	// Load the client krb5 config
	conf, err := config.NewConfigFromString(string(kdc))
	if err != nil {
		l.Fatalf("could not load krb5.conf: %v", err)
	}

	// Create the client with the keytab
	cl := client.NewClientWithKeytab("client/admin", "EXAMPLE.ORG", kt, conf, client.Logger(l), client.DisablePAFXFAST(true))

	// Log in the client
	err = cl.Login()
	if err != nil {
		l.Fatalf("could not login client: %v", err)
	}

	// Form the request
	url := "http://localhost" + port
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		l.Fatalf("could create request: %v", err)
	}

	spnegoCl := spnego.NewClient(cl, nil, "server/admin")

	// Make the request
	resp, err := spnegoCl.Do(r)
	if err != nil {
		l.Fatalf("error making request: %v", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Fatalf("error reading response body: %v", err)
	}
	fmt.Println(string(b))
}
