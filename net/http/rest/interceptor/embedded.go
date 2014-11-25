// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/shelter/protocol"
	"gopkg.in/mgo.v2"
	"net"
	"reflect"
)

type JSONCompliant struct {
	request  reflect.Value
	response reflect.Value
	message  protocol.Translator
}

func (j *JSONCompliant) RequestValue() reflect.Value {
	return j.request
}

func (j *JSONCompliant) SetRequestValue(r reflect.Value) {
	j.request = r
}

func (j *JSONCompliant) ResponseValue() reflect.Value {
	return j.response
}

func (j *JSONCompliant) SetResponseValue(r reflect.Value) {
	j.response = r
}

func (j *JSONCompliant) Message() protocol.Translator {
	return j.message
}

func (j *JSONCompliant) SetMessage(m protocol.Translator) {
	j.message = m
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

type RemoteAddressCompliant struct {
	remoteAddress net.IP
}

func (r *RemoteAddressCompliant) SetRemoteAddress(a net.IP) {
	r.remoteAddress = a
}

func (r *RemoteAddressCompliant) RemoteAddress() net.IP {
	return r.remoteAddress
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

type DatabaseCompliant struct {
	session  *mgo.Session  // MongoDB session
	database *mgo.Database // Database connection of the MongoDB session
}

func (d *DatabaseCompliant) SetDBSession(s *mgo.Session) {
	d.session = s
}

func (d *DatabaseCompliant) DBSession() *mgo.Session {
	return d.session
}

func (d *DatabaseCompliant) SetDB(database *mgo.Database) {
	d.database = database
}

func (d *DatabaseCompliant) DB() *mgo.Database {
	return d.database
}
