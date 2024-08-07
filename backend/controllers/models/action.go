package models

import (
	"errors"
	"github.com/linxlib/conv"
	"strconv"
)

// ActionType
// @Enum
type ActionType int

func (a *ActionType) MarshalJSON() ([]byte, error) {
	switch *a {
	case Reboot:
		return conv.Bytes("\"Reboot\""), nil
	case Stop:
		return conv.Bytes("\"Stop\""), nil
	case Start:
		return conv.Bytes(`"Start"`), nil
	default:
		return nil, errors.New("unknown action type")
	}
}

func (a *ActionType) UnmarshalJSON(bytes []byte) error {

	a1 := conv.String(bytes)
	a1, _ = strconv.Unquote(a1)
	switch a1 {
	case "Reboot":
		*a = Reboot
	case "Stop":
		*a = Stop
	case "Start":
		*a = Start
	default:
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
