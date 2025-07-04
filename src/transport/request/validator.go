package request

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(s any) error {	
	if err := cv.Validator.Struct(s); err != nil {
    return err
  }
  return nil
}