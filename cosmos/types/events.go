package types

import (
	"fmt"
	"strings"

	"github.com/functionx/go-sdk/tendermint/abci"
)

type Event abci.Event

type Events []Event

func (a *Attribute) String() string {
	return fmt.Sprintf("%s: %s", a.Key, a.Value)
}

type StringEvents []StringEvent

func (e StringEvents) String() string {
	var sb strings.Builder
	for _, e := range e {
		sb.WriteString(fmt.Sprintf("\t\t- %s\n", e.Type))
		for _, attr := range e.Attributes {
			sb.WriteString(fmt.Sprintf("\t\t\t- %s\n", attr.String()))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
