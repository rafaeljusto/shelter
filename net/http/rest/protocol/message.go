// protocol - REST protocol description
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package protocol

// MessageResponse struct was created to return a message to the user with more
// information to easy integrate and solve problems
type MessageResponse struct {
	Id      string `json:"id,omitempty"`      // Code for integration systems to automatically solve the problem
	Message string `json:"message,omitempty"` // Message in the user's desired language
	Links   []Link `json:"links,omitempty"`   // Links associating this message with other resources
}
