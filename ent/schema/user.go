package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("username").
			NotEmpty().
			Unique(),
		field.String("password").
			NotEmpty(),
		field.String("email"). // New field
					Optional().
					Unique(),
	}
}
