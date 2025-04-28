package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Item struct {
	ent.Schema
}

func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.Int("price"),
		field.String("description").
			Optional(),
	}
}
