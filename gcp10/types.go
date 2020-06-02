
package gcp10

import (
	"fmt"
)

type State struct {
	From   string
	To     string
	Amount string
}

type Gcp10TransferEvent struct {
	Name   string
	From   string
	To     string
	Amount *string
}

func (this *Gcp10TransferEvent) String() string {
	return fmt.Sprintf("name %s, from %s, to %s, amount %s", this.Name, this.From, this.To,
		this.Amount)
}
