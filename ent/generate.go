package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./schema
//go:generate go run -mod=mod entgo.io/contrib/entproto/cmd/entproto -path ./schema
