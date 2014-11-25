// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"encoding/json"
	"github.com/rafaeljusto/shelter/messages"
	"github.com/rafaeljusto/shelter/protocol"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type requestResponser interface {
	RequestValue() reflect.Value
	SetRequestValue(reflect.Value)
	ResponseValue() reflect.Value
	SetResponseValue(reflect.Value)
	Message() protocol.Translator
	SetMessage(protocol.Translator)
	Language() *messages.LanguagePack
}

type JSON struct {
	handler requestResponser
}

func NewJSON(h requestResponser) *JSON {
	return &JSON{handler: h}
}

func (i *JSON) Before(w http.ResponseWriter, r *http.Request) {
	m := strings.ToLower(r.Method)
	i.parse(m)

	if request := i.handler.RequestValue(); request.IsValid() {
		decoder := json.NewDecoder(r.Body)
		for {
			if err := decoder.Decode(request.Addr().Interface()); err == io.EOF {
				break
			} else if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				i.handler.SetMessage(protocol.NewMessageResponse(protocol.ErrorCodeInvalidJSONContent, nil))

				return
			}
		}
	}
}

func (i *JSON) After(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if !r.Close && r.Body != nil {
			err := r.Body.Close()

			if err != nil {
				// TODO!
				//i.handler.Logger().Error(err)
			}
		}
	}()

	if message := i.handler.Message(); message != nil {
		if message.Translate(i.handler.Language()) {
			if body, err := json.Marshal(message); err == nil {
				w.Write(body)

			} else {
				// TODO!
				//i.handler.Logger().Errorf("Error writing response. Details: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}

		return
	}

	response := i.handler.ResponseValue()

	if !response.IsValid() {
		return
	}

	elem := response.Interface()
	elemType := reflect.TypeOf(elem)

	// We are also checking for map types because they work like pointers
	if (elemType.Kind() == reflect.Ptr ||
		elemType.Kind() == reflect.Map ||
		elemType.Kind() == reflect.Slice) &&
		response.IsNil() {

		return
	}

	body, err := json.Marshal(elem)

	if err != nil {
		//i.handler.Logger().Errorf("Error writing response. Details: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// For HTTP HEAD method we never add the body, but we add the other headers as a
	// normal GET method. For more information check the RFC 2616 - 14.13.
	if len(body) > 0 && r.Method != "HEAD" {
		w.Write(body)
	}
}

func (i *JSON) parse(m string) {
	st := reflect.ValueOf(i.handler).Elem()

	for j := 0; j < st.NumField(); j++ {
		field := st.Type().Field(j)

		value := field.Tag.Get("request")
		if value == "all" || strings.Contains(value, m) {
			i.handler.SetRequestValue(st.Field(j))
			continue
		}

		value = field.Tag.Get("response")
		if value == "all" || strings.Contains(value, m) {
			i.handler.SetResponseValue(st.Field(j))
		}
	}
}
