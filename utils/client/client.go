package main

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
)

type shelterRequest struct {
	url    string
	method string
	body   string
	secret string
}

func main() {
	app := cli.NewApp()
	app.Name = "client"
	app.Usage = "Easy way to send requests to the REST server"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url,u",
			Usage: "URL to send the request",
		},
		cli.StringFlag{
			Name:  "method,m",
			Value: "GET",
			Usage: "HTTP method",
		},
		cli.StringFlag{
			Name:  "body,b",
			Usage: "HTTP body message",
		},
		cli.StringFlag{
			Name:  "file,f",
			Usage: "File to send as HTTP body message",
		},
		cli.StringFlag{
			Name:  "secret,s",
			Value: "abc123",
			Usage: "Secret used between client and server",
		},
	}

	app.Action = func(c *cli.Context) {
		shelterRequest := shelterRequest{
			url:    c.String("url"),
			method: c.String("method"),
			body:   c.String("body"),
			secret: c.String("secret"),
		}

		if shelterRequest.url == "" {
			fmt.Println("URL not informed")
			return
		}

		if filename := c.String("file"); filename != "" {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Println("Error opening file. Details:", err)
				return
			}

			content, err := ioutil.ReadAll(file)
			file.Close()

			if err != nil {
				fmt.Println("Error reading file. Details:", err)
				return
			}

			shelterRequest.body = string(content)
		}

		request, response, err := send(shelterRequest)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()

		if data, err := httputil.DumpRequestOut(request, true); err != nil {
			fmt.Println("Error dumping request. Details:", err)

		} else {
			printHTTP("REQUEST", string(data))
		}

		if data, err := httputil.DumpResponse(response, true); err != nil {
			fmt.Println("Error dumping response. Details:", err)
		} else {
			printHTTP("RESPONSE", string(data))
		}
	}

	app.Run(os.Args)
}

func send(shelterRequest shelterRequest) (*http.Request, *http.Response, error) {
	var request *http.Request
	var err error

	if len(shelterRequest.body) > 0 {
		if request, err = http.NewRequest(shelterRequest.method, shelterRequest.url,
			strings.NewReader(shelterRequest.body)); err != nil {

			return request, nil, fmt.Errorf("Error creating request. Details: %s", err)
		}

		hash := md5.New()
		hash.Write([]byte(shelterRequest.body))
		hashBytes := hash.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

		request.Header.Add("Content-MD5", hashBase64)

	} else {
		if request, err = http.NewRequest(shelterRequest.method, shelterRequest.url, nil); err != nil {
			return request, nil, fmt.Errorf("Error creating request. Details: %s", err)
		}
	}

	request.Header.Add("Accept", check.SupportedContentType)
	request.Header.Add("Accept-Charset", check.SupportedCharset)
	request.Header.Add("Accept-Language", "en-us")
	request.Header.Add("Content-Type", check.SupportedContentType+"; charset="+check.SupportedCharset)
	request.Header.Add("Date", time.Now().Format(time.RFC1123))

	stringToSign, err := check.BuildStringToSign(request, "key01")
	if err != nil {
		return request, nil, fmt.Errorf("Error creating authorization. Details: %s", err)
	}

	authorization := check.GenerateSignature(stringToSign, shelterRequest.secret)
	request.Header.Add("Authorization", fmt.Sprintf("%s %s:%s", check.SupportedNamespace, "key01", authorization))

	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}

	client := http.Client{
		Transport: transport,
	}

	response, err := client.Do(request)
	if err != nil {
		return request, nil, fmt.Errorf("Error sending request. Details: %s", err)
	}

	// Add again the body for the dump
	request.Body = ioutil.NopCloser(strings.NewReader(shelterRequest.body))
	return request, response, nil
}

func printHTTP(label, data string) {
	headersBodyIndex := strings.Index(string(data), "\r\n\r\n")
	if headersBodyIndex < 0 {
		fmt.Printf("[%s]\n\n%s\n\n", label, data)
		return
	}

	headers := data[:headersBodyIndex]
	headers = strings.TrimSpace(headers)

	body := data[headersBodyIndex+1:]
	body = strings.TrimSpace(body)

	if body == "" {
		fmt.Printf("[%s]\n\n%s\n\n", label, data)
		return
	}

	var jsonBody bytes.Buffer
	if err := json.Indent(&jsonBody, []byte(body), "", "  "); err != nil {
		fmt.Println("Error parsing JSON. Details:", err)
	}

	fmt.Printf("[%s]\n\n%s\n\n%s\n", label, headers, jsonBody.String())
}
