package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("nickname").
			Default("unknown"),
		field.Int("age").
			Positive(),
		field.Float("height").
			Positive().
			Default(0),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
