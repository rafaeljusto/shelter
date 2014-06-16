// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type CodecHandler interface {
	GetLanguage() *messages.LanguagePack
	MessageResponse(string, string) error
}

type JSONCodec struct {
	codecHandler CodecHandler
	errPosition  int
	reqPosition  int
	resPosition  int
}

func NewJSONCodec(h CodecHandler) *JSONCodec {
	return &JSONCodec{codecHandler: h}
}

func (c *JSONCodec) Before(w http.ResponseWriter, r *http.Request) {
	m := strings.ToLower(r.Method)
	c.parse(m)

	if c.reqPosition >= 0 {
		st := reflect.ValueOf(c.codecHandler).Elem()
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(st.Field(c.reqPosition).Addr().Interface()); err != nil {
			log.Println("Received an invalid JSON. Details:", err)

			if err := c.codecHandler.MessageResponse("invalid-json-content", r.URL.RequestURI()); err == nil {
				w.WriteHeader(http.StatusBadRequest)

			} else {
				log.Println("Error while writing response. Details:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}
	}
}

func (c *JSONCodec) After(w http.ResponseWriter, r *http.Request) {
	st := reflect.ValueOf(c.codecHandler).Elem()

	// We are going to always send the content HTTP header fields even if we don't have a
	// content, because if we don't set the GoLang HTTP server will add "text/plain"
	w.Header().Set("Content-Type", fmt.Sprintf("application/vnd.shelter+json; charset=%s", check.SupportedCharset))

	if c.errPosition >= 0 {
		if elem := st.Field(c.errPosition).Interface(); elem != nil {
			elemType := reflect.TypeOf(elem)

			if elemType.Kind() == reflect.Ptr && !st.Field(c.errPosition).IsNil() {
				body, err := json.Marshal(elem)
				if err != nil {
					log.Println("Error writing message response. Details:", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				c.addResponseHeaders(w, body)
				w.Write(body)
				return
			}
		}
	}

	if c.resPosition >= 0 {
		elem := st.Field(c.resPosition).Interface()
		elemType := reflect.TypeOf(elem)
		if elemType.Kind() == reflect.Ptr && st.Field(c.resPosition).IsNil() {
			return
		}

		body, err := json.Marshal(elem)
		if err != nil {
			log.Println("Error writing message response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.addResponseHeaders(w, body)

		// For HTTP HEAD method we never add the body, but we add the other headers as a
		// normal GET method. For more information check the RFC 2616 - 14.13.
		if len(body) > 0 && r.Method != "HEAD" {
			w.Write(body)
		}

	} else {
		c.addResponseHeaders(w, nil)
	}
}

func (c *JSONCodec) parse(m string) {
	c.errPosition, c.reqPosition, c.resPosition = -1, -1, -1

	st := reflect.ValueOf(c.codecHandler).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		if field.Tag == "error" {
			c.errPosition = i
			continue
		}

		value := field.Tag.Get("request")
		if value == "all" || strings.Contains(value, m) {
			c.reqPosition = i
			continue
		}

		value = field.Tag.Get("response")
		if value == "all" || strings.Contains(value, m) {
			c.resPosition = i
		}
	}
}

func (c *JSONCodec) addResponseHeaders(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))

	if l := c.codecHandler.GetLanguage(); l != nil {
		w.Header().Set("Content-Language", l.Name())
	}

	if len(body) > 0 {
		hash := md5.New()
		hash.Write(body)
		hashBytes := hash.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
		w.Header().Set("Content-MD5", hashBase64)
	}

	w.Header().Set("Accept", check.SupportedContentType)
	w.Header().Set("Accept-Language", messages.ShelterRESTLanguagePacks.Names())
	w.Header().Set("Accept-Charset", "utf-8")
	w.Header().Set("Accept-Encoding", "gzip")
	w.Header().Set("Date", time.Now().UTC().Format(time.RFC1123))
}
