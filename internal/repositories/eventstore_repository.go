package repositories

import (
	"context"

	"github.com/longgggwww/hrm-ms-permission/internal/ent"
)

// EventStoreRepository defines the interface for EventStore operations.
type EventStoreRepository interface {
	Create(ctx context.Context, event *ent.EventStore) (*ent.EventStore, error)
	GetByID(ctx context.Context, id int) (*ent.EventStore, error)
	Update(ctx context.Context, event *ent.EventStore) (*ent.EventStore, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*ent.EventStore, error)
}

// eventStoreRepository is the implementation of EventStoreRepository.
type eventStoreRepository struct {
	client *ent.Client
}

// NewEventStoreRepository creates a new instance of EventStoreRepository.
func NewEventStoreRepository(client *ent.Client) EventStoreRepository {
	return &eventStoreRepository{client: client}
}

func (r *eventStoreRepository) Create(ctx context.Context, event *ent.EventStore) (*ent.EventStore, error) {
	return r.client.EventStore.Create().
		SetEventID(event.EventID).
		SetAggregateID(event.AggregateID).
		SetEventType(event.EventType).
		SetPayload(event.Payload).
		SetCreatedAt(event.CreatedAt).
		Save(ctx)
}

func (r *eventStoreRepository) GetByID(ctx context.Context, id int) (*ent.EventStore, error) {
	return r.client.EventStore.Get(ctx, id)
}

func (r *eventStoreRepository) Update(ctx context.Context, event *ent.EventStore) (*ent.EventStore, error) {
	return r.client.EventStore.UpdateOne(event).
		SetEventID(event.EventID).
		SetAggregateID(event.AggregateID).
		SetEventType(event.EventType).
		SetPayload(event.Payload).
		SetCreatedAt(event.CreatedAt).
		Save(ctx)
}

func (r *eventStoreRepository) Delete(ctx context.Context, id int) error {
	return r.client.EventStore.DeleteOneID(id).Exec(ctx)
}

func (r *eventStoreRepository) List(ctx context.Context, limit, offset int) ([]*ent.EventStore, error) {
	return r.client.EventStore.Query().
		Limit(limit).
		Offset(offset).
		All(ctx)
}
