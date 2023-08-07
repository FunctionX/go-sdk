package tool

import (
	"reflect"

	"github.com/pkg/errors"
)

type Validate interface {
	ValidateBase() error
}

func ValidateParams(v Validate) error {
	if v == nil || reflect.ValueOf(v).IsNil() {
		return errors.New("request params is nil")
	}
	if err := v.ValidateBase(); err != nil {
		return err
	}
	return nil
}
