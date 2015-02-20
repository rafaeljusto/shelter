// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package mongodb is an interface to initialize a mgo connection
package mongodb

import (
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/rafaeljusto/shelter/secret"
	"sync"
	"time"
)

var (
	// List of functions that every DAO will register to add the necessary indexes that
	// the application must have to garantee the performance
	indexFunctions []func(*mgo.Database) error

	// MongoDB session singleton. According to Gustavo Niemeyer, created of mgo, we should
	// create only one session for the application and copy this object for every request.
	// The mgo library will be responsable for pooling the connections. For more information
	// check the threads https://groups.google.com/forum/#!topic/mgo-users/s1juysWHO8w and
	// https://groups.google.com/forum/#!topic/mgo-users/tDmuztwViDg
	session     *mgo.Session
	sessionLock sync.Mutex

	// Timeout is the amount of time to wait for a server to respond when
	// first connecting and on follow up operations in the session. If
	// timeout is zero, the call may block forever waiting for a connection
	// to be established.
	Timeout = time.Duration(5) * time.Second
)

// Open a new connection to a MongoDB database. For now we are using all default timeout
// values. The production database for the project will be shelter, but for tests
// purpouses we make this parameter configurable. We are also returning the session,
// because the caller must close it after use
func Open(
	uris []string,
	databaseName string,
	auth bool,
	username, password string,
) (*mgo.Database, *mgo.Session, error) {

	if err := initializeSession(uris, databaseName, auth, username, password); err != nil {
		return nil, nil, err
	}

	// Copying the session object as recomended by Gustavo Niemeyer
	localSession := session.Copy()

	// Choose the database
	database := localSession.DB(databaseName)

	// Apply all registered indexes
	for _, indexFunction := range indexFunctions {
		// If the database already have the index the function call will have no cost.
		// Depending on how the DAO add the index the operation can block until it ends or
		// not
		if err := indexFunction(database); err != nil {
			localSession.Close()
			return nil, nil, err
		}
	}

	return database, localSession, nil
}

// RegisterIndexFunction is the public function where the DAOs can register the indexes in
// their collections with the properties that they want
func RegisterIndexFunction(indexFunction func(*mgo.Database) error) {
	indexFunctions = append(indexFunctions, indexFunction)
}

// initializeSession verify if the session was already created, if not creates it
// otherwise use the already created session. We are reusing the session because Gustavo
// Niemeyer said that this should be the correct behaviour of the system
func initializeSession(
	uris []string,
	databaseName string,
	auth bool,
	username, password string,
) error {

	sessionLock.Lock()
	defer sessionLock.Unlock()

	if session != nil {
		return nil
	}

	var dialInfo mgo.DialInfo
	var err error

	if auth {
		password, err = secret.Decrypt(password)
		if err != nil {
			return err
		}

		dialInfo = mgo.DialInfo{
			Addrs:    uris,
			Timeout:  Timeout,
			Database: databaseName,
			Username: username,
			Password: password,
		}

	} else {
		dialInfo = mgo.DialInfo{
			Addrs:    uris,
			Timeout:  Timeout,
			Database: databaseName,
		}
	}

	// Connect to the database
	session, err = mgo.DialWithInfo(&dialInfo)
	return err
}
