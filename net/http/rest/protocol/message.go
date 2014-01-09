package protocol

// MessageResponse struct was created to return a message to the user with more
// information to easy integrate and solve problems
type MessageResponse struct {
	Id      string `json:"id,omitempty"`      // Code for integration systems to automatically solve the problem
	Message string `json:"message,omitempty"` // Message in the user's desired language
	Links   Links  `json:"links,omitempty"`   // Links associating this message with other resources
}
