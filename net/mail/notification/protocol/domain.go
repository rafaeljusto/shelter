package protocol

import (
	"github.com/rafaeljusto/shelter/model"
)

// Struct created to add the extra information necessary to build an e-mail template that is going
// to be used to notify the domain's owners
type Domain struct {
	model.Domain          // Domain object
	From         string   // E-mails from header
	To           []string // List of owner's e-mails to be alerted
}
