package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// UserPerm holds the schema definition for the UserPerm entity.
type UserPerm struct {
	ent.Schema
}

// Fields of the UserPerm.
func (UserPerm) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").NotEmpty().Annotations(entproto.Field(2)),
		field.UUID("perm_id", uuid.UUID{}).Annotations(entproto.Field(3)),
	}
}

// Edges of the UserPerm.
func (UserPerm) Edges() []ent.Edge {
	return nil
}

// Indexes of the UserPerm.
func (UserPerm) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("perm_id", "user_id").Unique(),
	}
}

// Annotations of the UserPerm.
func (UserPerm) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
