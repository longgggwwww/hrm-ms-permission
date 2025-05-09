// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/longgggwwww/hrm-ms-permission/ent/permgroup"
)

// PermGroup is the model entity for the PermGroup schema.
type PermGroup struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Code holds the value of the "code" field.
	Code string `json:"code,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PermGroupQuery when eager-loading is set.
	Edges        PermGroupEdges `json:"edges"`
	selectValues sql.SelectValues
}

// PermGroupEdges holds the relations/edges for other nodes in the graph.
type PermGroupEdges struct {
	// Perms holds the value of the perms edge.
	Perms []*Perm `json:"perms,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// PermsOrErr returns the Perms value or an error if the edge
// was not loaded in eager-loading.
func (e PermGroupEdges) PermsOrErr() ([]*Perm, error) {
	if e.loadedTypes[0] {
		return e.Perms, nil
	}
	return nil, &NotLoadedError{edge: "perms"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PermGroup) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case permgroup.FieldCode, permgroup.FieldName:
			values[i] = new(sql.NullString)
		case permgroup.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PermGroup fields.
func (pg *PermGroup) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case permgroup.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pg.ID = *value
			}
		case permgroup.FieldCode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field code", values[i])
			} else if value.Valid {
				pg.Code = value.String
			}
		case permgroup.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				pg.Name = value.String
			}
		default:
			pg.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the PermGroup.
// This includes values selected through modifiers, order, etc.
func (pg *PermGroup) Value(name string) (ent.Value, error) {
	return pg.selectValues.Get(name)
}

// QueryPerms queries the "perms" edge of the PermGroup entity.
func (pg *PermGroup) QueryPerms() *PermQuery {
	return NewPermGroupClient(pg.config).QueryPerms(pg)
}

// Update returns a builder for updating this PermGroup.
// Note that you need to call PermGroup.Unwrap() before calling this method if this PermGroup
// was returned from a transaction, and the transaction was committed or rolled back.
func (pg *PermGroup) Update() *PermGroupUpdateOne {
	return NewPermGroupClient(pg.config).UpdateOne(pg)
}

// Unwrap unwraps the PermGroup entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pg *PermGroup) Unwrap() *PermGroup {
	_tx, ok := pg.config.driver.(*txDriver)
	if !ok {
		panic("ent: PermGroup is not a transactional entity")
	}
	pg.config.driver = _tx.drv
	return pg
}

// String implements the fmt.Stringer.
func (pg *PermGroup) String() string {
	var builder strings.Builder
	builder.WriteString("PermGroup(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pg.ID))
	builder.WriteString("code=")
	builder.WriteString(pg.Code)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(pg.Name)
	builder.WriteByte(')')
	return builder.String()
}

// PermGroups is a parsable slice of PermGroup.
type PermGroups []*PermGroup
