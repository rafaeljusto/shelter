package handler

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"net/http"
	"time"
)

// List of errors that can occur while using the functions from this file. Low level
// errors can also be thrown
var (
	// To send the requests to a REST server we must have at least one address defined in
	// the configuration file
	ErrNoRESTAddresses = errors.New("Don't known the address of the REST server")

	// The user should configure a secret to sign the requests sent to the REST server.
	// Otherwise this error will be thrown on every request sent attempt
	ErrNoSecretFound = errors.New("No secret found to sign the REST request")
)

var (
	client http.Client
)

func init() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = http.Client{
		Transport: transport,
	}
}

func retrieveRESTAddress() (string, error) {
	var restAddresses []string
	for _, ln := range config.ShelterConfig.RESTServer.Listeners {
		if ln.TLS {
			restAddresses = append(restAddresses, fmt.Sprintf("https://[%s]:%d", ln.IP, ln.Port))
		} else {
			restAddresses = append(restAddresses, fmt.Sprintf("http://[%s]:%d", ln.IP, ln.Port))
		}
	}

	if len(restAddresses) == 0 {
		return "", ErrNoRESTAddresses
	}

	// Use the first REST address found in configuration file
	return restAddresses[0], nil
}

func signAndSend(r *http.Request, content []byte) (*http.Response, error) {
	r.Header.Set("Date", time.Now().Format(time.RFC1123))

	if r.ContentLength > 0 && content != nil {
		r.Header.Set("Content-Type", check.SupportedContentType)

		hash := md5.New()
		hash.Write(content)
		hashBytes := hash.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

		r.Header.Set("Content-MD5", hashBase64)
	}

	var key, secret string
	for key, secret = range config.ShelterConfig.RESTServer.Secrets {
		break
	}

	if len(key) == 0 || len(secret) == 0 {
		return nil, ErrNoSecretFound
	}

	stringToSign, err := check.BuildStringToSign(r, key)
	if err != nil {
		return nil, err
	}

	signature := check.GenerateSignature(stringToSign, secret)
	r.Header.Set("Authorization",
		fmt.Sprintf("%s %s:%s", check.SupportedNamespace, key, signature))

	return client.Do(r)
}
