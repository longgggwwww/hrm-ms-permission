package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PermGroup holds the schema definition for the PermGroup entity.
type PermGroup struct {
	ent.Schema
}

// Fields of the PermGroup.
func (PermGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").Unique().Annotations(entproto.Field(2)),
		field.String("name").Annotations(entproto.Field(3)),
	}
}

// Edges of the PermGroup.
func (PermGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("perms", Perm.Type).Ref("group").Annotations(entproto.Field(4)),
	}
}

// Method annotations for PermGroup entity
func (PermGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
