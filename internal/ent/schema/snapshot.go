package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Snapshot holds the schema definition for the Snapshot entity.
type Snapshot struct {
	ent.Schema
}

// Fields of the Snapshot.
func (Snapshot) Fields() []ent.Field {
	return []ent.Field{
		field.String("snapshot_id").Unique(),
		field.String("aggregate_id"),
		field.JSON("state", map[string]interface{}{}),
		field.Time("created_at"),
	}
}

// Edges of the Snapshot.
func (Snapshot) Edges() []ent.Edge {
	return nil
}

// Indexes of the Snapshot.
func (Snapshot) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("aggregate_id").
			Unique(),
	}
}
