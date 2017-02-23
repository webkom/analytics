package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/lib/pq"
	"time"
)

type Context map[string]interface{}

func (c *Context) Marshal() ([]byte, error) {
	jsonb, err := json.Marshal(c)
	return jsonb, err
}

func (c *Context) Value() (driver.Value, error) {
	if c == nil {
		return []byte("null"), nil
	}

	jsonb, err := c.Marshal()

	return string(jsonb), err
}

type Event struct {
	AnonymousId string    `json:"anonymousId"`
	UserId      string    `json:"userId"`
	Type        string    `json:"type"`
	Context     *Context  `json:"context"`
	ReceivedAt  time.Time `json:"receivedAt"`
	SentAt      time.Time `json:"sentAt"`
	Timestamp   time.Time `json:"timestamp"`
}

type BatchEvents struct {
	Batch  []*Event  `json:"batch"`
	SentAt time.Time `json:"sentAt"`
}

func (be *BatchEvents) createBatchEvents(db *sql.DB) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"events", "anonymous_id", "user_id",
		"type", "context", "received_at", "sent_at", "timestamp",
	))
	if err != nil {
		return err
	}

	receivedAt := time.Now().UTC()

	for _, event := range be.Batch {

		_, err = stmt.Exec(
			event.AnonymousId, event.UserId, event.Type, event.Context,
			receivedAt, be.SentAt, event.Timestamp,
		)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}
