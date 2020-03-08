package models

import (
	"bytes"
	"fmt"
	"time"

	"github.com/delivc/identity/storage"
	"github.com/delivc/identity/storage/namespace"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// AuditAction holds different action types as a string
type AuditAction string
type auditLogType string

const (
	// LoginAction is a action type
	LoginAction AuditAction = "login"
	// LogoutAction is a action type
	LogoutAction AuditAction = "logout"
	// InviteAcceptedAction is a action type
	InviteAcceptedAction AuditAction = "invite_accepted"
	// UserSignedUpAction is a action type
	UserSignedUpAction AuditAction = "user_signedup"
	// UserInvitedAction is a action type
	UserInvitedAction AuditAction = "user_invited"
	// UserDeletedAction is a action type
	UserDeletedAction AuditAction = "user_deleted"
	// UserModifiedAction is a action type
	UserModifiedAction AuditAction = "user_modified"
	// UserRecoveryRequestedAction is a action type
	UserRecoveryRequestedAction AuditAction = "user_recovery_requested"
	// TokenRevokedAction is a action type
	TokenRevokedAction AuditAction = "token_revoked"
	// TokenRefreshedAction is a action type
	TokenRefreshedAction AuditAction = "token_refreshed"

	account auditLogType = "account"
	team    auditLogType = "team"
	token   auditLogType = "token"
	user    auditLogType = "user"
)

var actionLogTypeMap = map[AuditAction]auditLogType{
	LoginAction:                 account,
	LogoutAction:                account,
	InviteAcceptedAction:        account,
	UserSignedUpAction:          team,
	UserInvitedAction:           team,
	UserDeletedAction:           team,
	TokenRevokedAction:          token,
	TokenRefreshedAction:        token,
	UserModifiedAction:          user,
	UserRecoveryRequestedAction: user,
}

// AuditLogEntry is the database model for audit log entries.
type AuditLogEntry struct {
	InstanceID uuid.UUID `json:"-" db:"instance_id"`
	ID         uuid.UUID `json:"id" db:"id"`

	Payload JSONMap `json:"payload" db:"payload"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TableName returns the audit tablename
func (AuditLogEntry) TableName() string {
	tableName := "audit_log_entries"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "_" + tableName
	}

	return tableName
}

// NewAuditLogEntry creates a new audit entry
func NewAuditLogEntry(tx *storage.Connection, instanceID uuid.UUID, actor *User, action AuditAction, traits map[string]interface{}) error {
	id, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "Error generating unique id")
	}
	l := AuditLogEntry{
		InstanceID: instanceID,
		ID:         id,
		Payload: JSONMap{
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"actor_id":    actor.ID,
			"actor_email": actor.Email,
			"action":      action,
			"log_type":    actionLogTypeMap[action],
		},
	}

	if name, ok := actor.UserMetaData["full_name"]; ok {
		l.Payload["actor_name"] = name
	}

	if traits != nil {
		l.Payload["traits"] = traits
	}

	return errors.Wrap(tx.Create(&l), "Database error creating audit log entry")
}

// FindAuditLogEntries searches for audit logs in the database
func FindAuditLogEntries(tx *storage.Connection, instanceID uuid.UUID, filterColumns []string, filterValue string, pageParams *Pagination) ([]*AuditLogEntry, error) {
	q := tx.Q().Order("created_at desc").Where("instance_id = ?", instanceID)

	if len(filterColumns) > 0 && filterValue != "" {
		lf := "%" + filterValue + "%"

		builder := bytes.NewBufferString("(")
		values := make([]interface{}, len(filterColumns))

		for idx, col := range filterColumns {
			builder.WriteString(fmt.Sprintf("payload->>'$.%s' COLLATE utf8mb4_unicode_ci LIKE ?", col))
			values[idx] = lf

			if idx+1 < len(filterColumns) {
				builder.WriteString(" OR ")
			}
		}
		builder.WriteString(")")

		q = q.Where(builder.String(), values...)
	}

	logs := []*AuditLogEntry{}
	var err error
	if pageParams != nil {
		err = q.Paginate(int(pageParams.Page), int(pageParams.PerPage)).All(&logs)
		pageParams.Count = uint64(q.Paginator.TotalEntriesSize)
	} else {
		err = q.All(&logs)
	}

	return logs, err
}
