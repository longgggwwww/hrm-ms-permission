package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty().Annotations(entproto.Field(2)),
		field.String("color").Optional().Annotations(entproto.Field(3)),
		field.String("description").Optional().Annotations(entproto.Field(4)),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("perms", Perm.Type).Annotations(entproto.Field(5)),
	}
}

// Method annotations for Role entity
func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
