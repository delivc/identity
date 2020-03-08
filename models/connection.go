package models

import (
	"github.com/delivc/identity/storage"
	"github.com/gobuffalo/pop"
)

// Pagination model
type Pagination struct {
	Page    uint64
	PerPage uint64
	Count   uint64
}

// Offset for pagination
func (p *Pagination) Offset() uint64 {
	return (p.Page - 1) * p.PerPage
}

// SortDirection holds Ascending or Descending cosnt
type SortDirection string

// Ascending sort direction
const Ascending SortDirection = "ASC"

// Descending sortdirection
const Descending SortDirection = "DESC"

// CreatedAt is a constant! ;O
const CreatedAt = "created_at"

// SortParams ?field,field,field
type SortParams struct {
	Fields []SortField
}

// SortField sort by what
type SortField struct {
	Name string
	Dir  SortDirection
}

// TruncateAll truncates all models
func TruncateAll(conn *storage.Connection) error {
	return conn.Transaction(func(tx *storage.Connection) error {
		if err := tx.RawQuery("TRUNCATE " + (&pop.Model{Value: User{}}).TableName()).Exec(); err != nil {
			return err
		}
		if err := tx.RawQuery("TRUNCATE " + (&pop.Model{Value: RefreshToken{}}).TableName()).Exec(); err != nil {
			return err
		}
		if err := tx.RawQuery("TRUNCATE " + (&pop.Model{Value: AuditLogEntry{}}).TableName()).Exec(); err != nil {
			return err
		}
		return tx.RawQuery("TRUNCATE " + (&pop.Model{Value: Instance{}}).TableName()).Exec()
	})
}
