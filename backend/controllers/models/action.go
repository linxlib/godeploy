package models

import (
	"errors"
	"github.com/linxlib/conv"
	"strconv"
)

// ActionType
// @Enum
type ActionType int

// MarshalJSON implements the json.Marshaler interface.
// It encodes the ActionType as a JSON string.
func (a *ActionType) MarshalJSON() ([]byte, error) {
	switch *a {
	case Reboot:
		return []byte(`"Reboot"`), nil
	case Stop:
		return []byte(`"Stop"`), nil
	case Start:
		return []byte(`"Start"`), nil
	default:
		return nil, errors.New("unknown action type")
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It decodes the ActionType from a JSON string.
func (a *ActionType) UnmarshalJSON(bytes []byte) error {
	// Remove quotes from the string
	a1 := conv.String(bytes)
	a1, _ = strconv.Unquote(a1)

	// Check if the ActionType is one of the known constants
	switch a1 {
	case "Reboot":
		*a = Reboot
	case "Stop":
		*a = Stop
	case "Start":
		*a = Start
	default:
		// If not, try to parse it as an integer
		*a = ActionType(conv.Int(conv.String(bytes)))
	}

	return nil
}

const (
	AnotherOne ActionType = 3
)

const (
	Reboot ActionType = iota //reboot
	Stop                     // stop
	Start                    // start
)
