package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Perm holds the schema definition for the Perm entity.
type Perm struct {
	ent.Schema
}

// Fields of the Perm.
func (Perm) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").Unique().Annotations(entproto.Field(2)),
		field.String("name").NotEmpty().Annotations(entproto.Field(3)),
		field.String("description").Optional().Annotations(entproto.Field(4)),
	}
}

// Edges of the Perm.
func (Perm) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("group", PermGroup.Type).Unique().Annotations(entproto.Field(5)),
		edge.From("roles", Role.Type).Ref("perms").Annotations(entproto.Field(6)),
	}
}

// Method annotations for Perm entity
func (Perm) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
