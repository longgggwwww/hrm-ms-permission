package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// EventStore holds the schema definition for the EventStore entity.
type EventStore struct {
	ent.Schema
}

// Fields of the EventStore.
func (EventStore) Fields() []ent.Field {
	return []ent.Field{
		field.String("event_id").Unique(),
		field.String("aggregate_id"),
		field.String("event_type"),
		field.JSON("payload", map[string]interface{}{}),
		field.Time("created_at"),
		field.Int("version").Default(1),
	}
}

// Edges of the EventStore.
func (EventStore) Edges() []ent.Edge {
	return nil
}

// Indexes of the EventStore.
func (EventStore) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("aggregate_id", "event_type").
			Unique(),
	}
}
