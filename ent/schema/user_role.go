package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// UserRole holds the schema definition for the UserRole entity.
type UserRole struct {
	ent.Schema
}

// Fields of the UserRole.
func (UserRole) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Annotations(entproto.Field(1)),
		field.String("user_id").NotEmpty().Annotations(entproto.Field(2)),
		field.UUID("role_id", uuid.UUID{}).Annotations(entproto.Field(3)),
		field.Time("created_at").Immutable().Default(time.Now).Annotations(entproto.Field(4)),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Annotations(entproto.Field(5)),
	}
}

// Edges of the UserRole.
func (UserRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).
			Ref("user_roles").
			Unique().
			Field("role_id").
			Required().
			Annotations(entproto.Field(6)),
	}
}

// Indexes of the UserRole.
func (UserRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "user_id").Unique(),
	}
}

// Annotations of the UserRole.
func (UserRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
