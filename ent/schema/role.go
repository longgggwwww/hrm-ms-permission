package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("code").
			Unique().
			NotEmpty().
			Annotations(entproto.Field(2)),
		field.String("name").
			NotEmpty().
			Annotations(entproto.Field(3)),
		field.String("color").
			Optional().
			Annotations(entproto.Field(4)),
		field.String("description").
			Optional().
			Annotations(entproto.Field(5)),
		field.Time("created_at").
			Nillable().
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(6)),
		field.Time("updated_at").
			Nillable().
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entproto.Field(7)),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("perms", Perm.Type).
			Annotations(entproto.Field(8)),
		edge.To("user_roles", UserRole.Type).
			Annotations(entproto.Field(9)),
	}
}

// Method annotations for Role entity
func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
