package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"time"
)

type Payload map[string]interface{}

func (p *Payload) Marshal() ([]byte, error) {
	jsonb, err := json.Marshal(p)
	return jsonb, err
}

func (p *Payload) Value() (driver.Value, error) {
	if p == nil {
		return []byte("null"), nil
	}

	jsonb, err := p.Marshal()
	return string(jsonb), err
}

type Identify struct {
	AnonymousId string    `json:"anonymousId"`
	UserId      string    `json:"userId"`
	Context     *Payload  `json:"context"`
	Traits      *Payload  `json:"traits"`
	ReceivedAt  time.Time `json:"receivedAt"`
	SentAt      time.Time `json:"sentAt"`
	Timestamp   time.Time `json:"timestamp"`
}

func (i *Identify) Save(txn *sql.Tx, be *BatchEvents) error {
	stmt, err := txn.Prepare(pq.CopyIn(
		"identifies", "anonymous_id", "user_id", "context", "traits", "received_at",
		"sent_at", "timestamp",
	))
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(
		i.AnonymousId, i.UserId, i.Context, i.Traits, be.ReceivedAt, i.SentAt, i.Timestamp,
	); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}

	return nil
}

type Page struct {
	AnonymousId string    `json:"anonymousId"`
	UserId      string    `json:"userId"`
	Context     *Payload  `json:"context"`
	Name        string    `json:"name"`
	Properties  *Payload  `json:"properties"`
	ReceivedAt  time.Time `json:"receivedAt"`
	SentAt      time.Time `json:"sentAt"`
	Timestamp   time.Time `json:"timestamp"`
}

func (p *Page) Save(txn *sql.Tx, be *BatchEvents) error {
	stmt, err := txn.Prepare(pq.CopyIn(
		"pages", "anonymous_id", "user_id", "context", "name", "properties",
		"received_at", "sent_at", "timestamp",
	))
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(
		p.AnonymousId, p.UserId, p.Context, p.Name, p.Properties, be.ReceivedAt, p.SentAt, p.Timestamp,
	); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}

	return nil
}

type Screen struct {
	AnonymousId string    `json:"anonymousId"`
	UserId      string    `json:"userId"`
	Context     *Payload  `json:"context"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Properties  *Payload  `json:"properties"`
	ReceivedAt  time.Time `json:"receivedAt"`
	SentAt      time.Time `json:"sentAt"`
	Timestamp   time.Time `json:"timestamp"`
}

func (s *Screen) Save(txn *sql.Tx, be *BatchEvents) error {
	stmt, err := txn.Prepare(pq.CopyIn(
		"screens", "anonymous_id", "user_id", "context", "name", "category", "properties",
		"received_at", "sent_at", "timestamp",
	))
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(
		s.AnonymousId, s.UserId, s.Context, s.Name, s.Category, s.Properties, be.ReceivedAt, s.SentAt,
		s.Timestamp,
	); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}

	return nil
}

type Track struct {
	AnonymousId string    `json:"anonymousId"`
	UserId      string    `json:"userId"`
	Context     *Payload  `json:"context"`
	Event       string    `json:"event"`
	Properties  *Payload  `json:"properties"`
	ReceivedAt  time.Time `json:"receivedAt"`
	SentAt      time.Time `json:"sentAt"`
	Timestamp   time.Time `json:"timestamp"`
}

func (t *Track) Save(txn *sql.Tx, be *BatchEvents) error {
	stmt, err := txn.Prepare(pq.CopyIn(
		"tracks", "anonymous_id", "user_id", "context", "event", "properties",
		"received_at", "sent_at", "timestamp",
	))
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(
		t.AnonymousId, t.UserId, t.Context, t.Event, t.Properties, be.ReceivedAt, t.SentAt, t.Timestamp,
	); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}

	return nil
}

type Group struct {
	AnonymousId string    `json:"anonymousId"`
	UserId      string    `json:"userId"`
	Context     *Payload  `json:"context"`
	Traits      *Payload  `json:"traits"`
	GroupId     string    `json:"group_id"`
	ReceivedAt  time.Time `json:"receivedAt"`
	SentAt      time.Time `json:"sentAt"`
	Timestamp   time.Time `json:"timestamp"`
}

func (g *Group) Save(txn *sql.Tx, be *BatchEvents) error {
	stmt, err := txn.Prepare(pq.CopyIn(
		"groups", "anonymous_id", "user_id", "context", "traits", "group_id", "received_at",
		"sent_at", "timestamp",
	))
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(
		g.AnonymousId, g.UserId, g.Context, g.Traits, g.GroupId, be.ReceivedAt, g.SentAt, g.Timestamp,
	); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}

	return nil
}

type BatchEvents struct {
	Batch      []Payload `json:"batch"`
	Context    Payload   `json:"context"`
	SentAt     time.Time `json:"sentAt"`
	ReceivedAt time.Time
}

func (be *BatchEvents) createBatchEvents(db *sql.DB) error {
	be.ReceivedAt = time.Now().UTC()

	txn, err := db.Begin()
	if err != nil {
		return err
	}

	for _, v := range be.Batch {

		switch v["type"] {
		case "identify":
			var identify Identify
			config := &mapstructure.DecoderConfig{
				DecodeHook: NormalizeTypesHookFunc(),
				Result:     &identify,
			}
			decoder, err := mapstructure.NewDecoder(config)
			if err != nil {
				return err
			}

			if err := decoder.Decode(v); err != nil {
				return err
			}
			if err := identify.Save(txn, be); err != nil {
				return err
			}
		case "track":
			var track Track
			config := &mapstructure.DecoderConfig{
				DecodeHook: NormalizeTypesHookFunc(),
				Result:     &track,
			}
			decoder, err := mapstructure.NewDecoder(config)
			if err != nil {
				return err
			}

			if err := decoder.Decode(v); err != nil {
				return err
			}
			if err := track.Save(txn, be); err != nil {
				return err
			}
		case "page":
			var page Page
			config := &mapstructure.DecoderConfig{
				DecodeHook: NormalizeTypesHookFunc(),
				Result:     &page,
			}
			decoder, err := mapstructure.NewDecoder(config)
			if err != nil {
				return err
			}

			if err := decoder.Decode(v); err != nil {
				return err
			}
			if err := page.Save(txn, be); err != nil {
				return err
			}
		case "screen":
			var screen Screen
			config := &mapstructure.DecoderConfig{
				DecodeHook: NormalizeTypesHookFunc(),
				Result:     &screen,
			}
			decoder, err := mapstructure.NewDecoder(config)
			if err != nil {
				return err
			}

			if err := decoder.Decode(v); err != nil {
				return err
			}
			if err := screen.Save(txn, be); err != nil {
				return err
			}
		case "group":
			var group Group
			config := &mapstructure.DecoderConfig{
				DecodeHook: NormalizeTypesHookFunc(),
				Result:     &group,
			}
			decoder, err := mapstructure.NewDecoder(config)
			if err != nil {
				return err
			}

			if err := decoder.Decode(v); err != nil {
				return err
			}
			if err := group.Save(txn, be); err != nil {
				return err
			}
		default:
			logrus.Warn("unknown_event_type")
		}

	}

	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}
