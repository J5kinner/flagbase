// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// ResourceKey Resource Key
//
// A unique key used to identify this resource.
//
// swagger:model ResourceKey
type ResourceKey string

// Validate validates this resource key
func (m ResourceKey) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validate.MinLength("", "body", string(m), 4); err != nil {
		return err
	}

	if err := validate.MaxLength("", "body", string(m), 30); err != nil {
		return err
	}

	if err := validate.Pattern("", "body", string(m), `^[a-z0-9]+([_ -]?[a-z0-9])*$`); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
