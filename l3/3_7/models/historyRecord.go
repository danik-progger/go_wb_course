package models

import "time"

// HistoryRecord represents a history record of changes
type HistoryRecord struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ItemID    uint      `json:"item_id" gorm:"index"`
	Action    string    `json:"action"` // INSERT, UPDATE, DELETE
	OldValues string    `json:"old_values,omitempty"`
	NewValues string    `json:"new_values,omitempty"`
	ChangedBy string    `json:"changed_by"`
	ChangedAt time.Time `json:"changed_at"`
}

// TableName specifies the table name for HistoryRecord
func (HistoryRecord) TableName() string {
	return "history_records"
}
