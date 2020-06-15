// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// SuccessResponse SuccessResponse
//
// Standard success response object
//
// swagger:model SuccessResponse
type SuccessResponse struct {

	// data
	Data interface{} `json:"data,omitempty"`
}

// Validate validates this success response
func (m *SuccessResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SuccessResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SuccessResponse) UnmarshalBinary(b []byte) error {
	var res SuccessResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
