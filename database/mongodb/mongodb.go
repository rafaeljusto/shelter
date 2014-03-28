// mongodb - Interface to initialize a mgo connection
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mongodb

import (
	"labix.org/v2/mgo"
	"sync"
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
)

// Open a new connection to a MongoDB database. For now we are using all default timeout
// values. The production database for the project will be shelter, but for tests
// purpouses we make this parameter configurable. We are also returning the session,
// because the caller must close it after use
func Open(uri, databaseName string) (*mgo.Database, *mgo.Session, error) {
	if err := initializeSession(uri); err != nil {
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
func initializeSession(uri string) error {
	sessionLock.Lock()
	defer sessionLock.Unlock()

	if session != nil {
		return nil
	}

	// Connect to the database
	var err error
	session, err = mgo.Dial(uri)
	return err
}
