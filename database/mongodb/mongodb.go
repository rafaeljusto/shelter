package mongodb

import (
	"labix.org/v2/mgo"
)

// Open a new connection to a MongoDB database. For now we are using all
// default timeout values. The production database for the project will be
// shelter, but for tests purpouses we make this parameter configurable
func Open(uri, database string) (*mgo.Database, error) {
	// Connect to the database
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}

	// Choose the database
	return session.DB(database), nil
}
