package util

import (
	"errors"
	"regexp"
)

type RegUtil struct {
}

func (*RegUtil) CanPraseBigFloat(str string) error {
	return canPraseBigFloat(str)
}

func canPraseBigFloat(str string) error {
	reg, err := regexp.Compile(`^[0-9]+\.{0,1}[0-9]*$`)
	if err != nil {
		return err
	}
	bo := reg.FindString(str)
	if bo == "" {
		return errors.New("String validation failed.")
	}
	return nil
}
