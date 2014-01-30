package mongodb

import (
	"labix.org/v2/mgo"
)

var (
	// List of functions that every DAO will register to add the necessary indexes that
	// the application must have to garantee the performance
	indexFunctions []func(*mgo.Database) error
)

// Open a new connection to a MongoDB database. For now we are using all default timeout
// values. The production database for the project will be shelter, but for tests
// purpouses we make this parameter configurable. We are also returning the session,
// because the caller must close it after use
func Open(uri, databaseName string) (*mgo.Database, *mgo.Session, error) {
	// Connect to the database
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, nil, err
	}

	// Choose the database
	database := session.DB(databaseName)

	// Apply all registered indexes
	for _, indexFunction := range indexFunctions {
		// If the database already have the index the function call will have no cost.
		// Depending on how the DAO add the index the operation can block until it ends or
		// not
		if err := indexFunction(database); err != nil {
			session.Close()
			return nil, nil, err
		}
	}

	return database, session, nil
}

// RegisterIndexFunction is the public function where the DAOs can register the indexes in
// their collections with the properties that they want
func RegisterIndexFunction(indexFunction func(*mgo.Database) error) {
	indexFunctions = append(indexFunctions, indexFunction)
}
