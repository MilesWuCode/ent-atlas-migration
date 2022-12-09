// Code generated by ent, DO NOT EDIT.

package ent

import (
	"ent-atlas-migration/ent/schema"
	"ent-atlas-migration/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[0].Descriptor()
	// user.DefaultName holds the default value on creation for the name field.
	user.DefaultName = userDescName.Default.(string)
	// userDescAge is the schema descriptor for age field.
	userDescAge := userFields[1].Descriptor()
	// user.AgeValidator is a validator for the "age" field. It is called by the builders before save.
	user.AgeValidator = userDescAge.Validators[0].(func(int) error)
	// userDescHeight is the schema descriptor for height field.
	userDescHeight := userFields[2].Descriptor()
	// user.DefaultHeight holds the default value on creation for the height field.
	user.DefaultHeight = userDescHeight.Default.(float64)
	// user.HeightValidator is a validator for the "height" field. It is called by the builders before save.
	user.HeightValidator = userDescHeight.Validators[0].(func(float64) error)
}
