// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/messages"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"github.com/rafaeljusto/shelter/protocol"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	shelterContentType = "application/vnd.shelter+json"
	shelterCharset     = "utf-8"
	timeFrameDuration  = time.Duration(10 * time.Minute)
)

type httpHeaders interface {
	SetMessage(protocol.Translator)
	Language() *messages.LanguagePack
}

type HTTPHeaders struct {
	handler httpHeaders
}

func NewHTTPHeaders(h httpHeaders) *HTTPHeaders {
	return &HTTPHeaders{handler: h}
}

func (i HTTPHeaders) Before(w http.ResponseWriter, r *http.Request) {
	var messagesHolder protocol.MessagesHolder
	if ok, err := check.HTTPDate(r); !ok || err != nil {
		messagesHolder.Add(protocol.NewMessageResponseWithField(protocol.ErrorCodeInvalidHeaderDate,
			"Date", r.Header.Get("Date"), nil))
	}
	if !check.HTTPContentType(r) {
		messagesHolder.Add(protocol.NewMessageResponseWithField(protocol.ErrorCodeInvalidContentType,
			"Content-Type", r.Header.Get("Content-Type"), nil))
	}

	var bodyBuffer bytes.Buffer
	io.Copy(&bodyBuffer, r.Body)
	r.Body = ioutil.NopCloser(&bodyBuffer)
	body, err := ioutil.ReadAll(&bodyBuffer)
	if err != nil {
		// TODO
	}

	if !check.HTTPContentMD5(r, body) {
		messagesHolder.Add(protocol.NewMessageResponseWithField(protocol.ErrorCodeInvalidContentMD5,
			"Content-MD5", r.Header.Get("Content-MD5"), nil))
	}

	// messagesHolder.Add(i.checkContentMD5(r))
	// messagesHolder.Add(i.checkAccept(r))
	// messagesHolder.Add(i.checkAcceptLanguage(r))
	// messagesHolder.Add(i.checkAcceptCharset(r))

	if messages := messagesHolder.Messages(); messages != nil {
		i.handler.SetMessage(messages)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (i HTTPHeaders) After(w http.ResponseWriter, r *http.Request) {
	responseWriter, ok := w.(*handy.BufferedResponseWriter)
	if !ok {
		//i.handler.Logger().Error(fmt.Errorf("Writer is not a handy framework object"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := responseWriter.Body.Bytes()

	// We are going to always send the content HTTP header fields even if we don't have a
	// content, because if we don't set the GoLang HTTP server will add "text/plain"
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", check.SupportedContentType, "utf-8"))
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))

	if i.handler.Language() != nil {
		w.Header().Set("Content-Language", i.handler.Language().Name())
	}

	if len(body) > 0 {
		hash := md5.New()
		hash.Write(body)
		contentMD5 := base64.StdEncoding.EncodeToString(hash.Sum(nil))
		w.Header().Set("Content-MD5", contentMD5)
	}

	w.Header().Set("Accept", check.SupportedContentType)
	w.Header().Set("Accept-Charset", shelterCharset)
	w.Header().Set("Accept-Encoding", "gzip")
	w.Header().Set("Accept-Language", messages.ShelterRESTLanguagePacks.Names())
	w.Header().Add("Date", time.Now().UTC().Format(time.RFC1123))
}
