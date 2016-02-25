package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
)

var (
	goServerCAFile = "go-server-ca.pem"
	privateKeyFile = "agent-private-key.pem"
	certFile       = "agent-cert.pem"
	serverHost     = "localhost"
	sslPort        = "8154"
	httpPort       = "8153"
)

func sslHostAndPort() string {
	return serverHost + ":" + sslPort
}

func httpsServerURL(path string) string {
	return "https://" + sslHostAndPort() + path
}

func httpServerURL(path string) string {
	return "http://" + serverHost + ":" + httpPort + path
}

func readGoServerCACert() {
	conn, err := tls.Dial("tcp", sslHostAndPort(), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	defer conn.Close()
	state := conn.ConnectionState()
	certOut, err := os.Create(goServerCAFile)
	if err != nil {
		log.Fatalf("failed to open %v for writing: %s", goServerCAFile, err)
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: state.PeerCertificates[0].Raw})
}

func main() {
	if serverHost == "localhost" {
		serverHost, _ = os.Hostname()
	}
	readGoServerCACert()
	register()
	caCert, err := ioutil.ReadFile(goServerCAFile)
	if err != nil {
		panic(err)
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(caCert))
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair(certFile, privateKeyFile)
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}
	resp, err := client.Get(httpsServerURL("/"))
	if err != nil {
		panic(err)
	}
	log.Println(resp.Body)
	defer resp.Body.Close()

}

type Registration struct {
	AgentPrivateKey, AgentCertificate string
}

func register() {
	hostname, _ := os.Hostname()
	workingDir, _ := os.Getwd()
	resp, _ := http.PostForm(httpServerURL("/go/admin/agent"),
		url.Values{
			"hostname":                      {hostname},
			"uuid":                          {"564e9408-fb78-4856-4215-52e0-e14bb054"},
			"location":                      {workingDir},
			"operatingSystem":               {runtime.GOOS},
			"usablespace":                   {"5000000000"},
			"agentAutoRegisterKey":          {""},
			"agentAutoRegisterResources":    {""},
			"agentAutoRegisterEnvironments": {""},
			"agentAutoRegisterHostname":     {""},
			"elasticAgentId":                {""},
			"elasticPluginId":               {""}})

	defer resp.Body.Close()
	var registration Registration
	dec := json.NewDecoder(resp.Body)

	err := dec.Decode(&registration)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(privateKeyFile, []byte(registration.AgentPrivateKey), 0600)
	ioutil.WriteFile(certFile, []byte(registration.AgentCertificate), 0600)
}

func createFile(name string) *os.File {
	file, err := os.Create(name)
	if err != nil {
		panic(fmt.Sprintf("failed to open %v for writing: %s", name, err))
	}
	return file
}
