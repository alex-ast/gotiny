// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetURLResponse get Url response
//
// swagger:model GetUrlResponse
type GetURLResponse struct {

	// source
	Source string `json:"source,omitempty"`

	// status
	Status *Status `json:"status,omitempty"`

	// url info
	URLInfo *URLInfo `json:"urlInfo,omitempty"`
}

// Validate validates this get Url response
func (m *GetURLResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateURLInfo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GetURLResponse) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if m.Status != nil {
		if err := m.Status.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("status")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("status")
			}
			return err
		}
	}

	return nil
}

func (m *GetURLResponse) validateURLInfo(formats strfmt.Registry) error {
	if swag.IsZero(m.URLInfo) { // not required
		return nil
	}

	if m.URLInfo != nil {
		if err := m.URLInfo.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("urlInfo")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("urlInfo")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get Url response based on the context it is used
func (m *GetURLResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateURLInfo(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GetURLResponse) contextValidateStatus(ctx context.Context, formats strfmt.Registry) error {

	if m.Status != nil {
		if err := m.Status.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("status")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("status")
			}
			return err
		}
	}

	return nil
}

func (m *GetURLResponse) contextValidateURLInfo(ctx context.Context, formats strfmt.Registry) error {

	if m.URLInfo != nil {
		if err := m.URLInfo.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("urlInfo")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("urlInfo")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *GetURLResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetURLResponse) UnmarshalBinary(b []byte) error {
	var res GetURLResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
