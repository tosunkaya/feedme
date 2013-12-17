package backend

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/zimmski/feedme"
)

type Postgresql struct {
	Db *sqlx.DB
}

func NewBackendPostgresql() Backend {
	return new(Postgresql)
}

func (p *Postgresql) Init(params BackendParameters) error {
	var err error

	p.Db, err = sqlx.Connect("postgres", params.Spec)

	if err != nil {
		return errors.New(fmt.Sprintf("Cannot connect to database: %v", err))
	}

	err = p.Db.Ping()

	if err != nil {
		return errors.New(fmt.Sprintf("Cannot ping database: %v", err))
	}

	p.Db.SetMaxIdleConns(params.MaxIdleConns)
	p.Db.SetMaxOpenConns(params.MaxOpenConns)

	return nil
}

func (p *Postgresql) FindFeed(feedName string) (*feedme.Feed, error) {
	feed := &feedme.Feed{}
	err := p.Db.Get(feed, "SELECT * FROM feeds WHERE name = $1", feedName)

	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	return feed, err
}

func (p *Postgresql) SearchFeeds() ([]feedme.Feed, error) {
	feeds := []feedme.Feed{}

	err := p.Db.Select(&feeds, "SELECT * FROM feeds ORDER BY name")

	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	return feeds, err
}

func (p *Postgresql) SearchItems(feed *feedme.Feed) ([]feedme.Item, error) {
	items := []feedme.Item{}
	err := p.Db.Select(&items, "SELECT * FROM items WHERE feed = $1 ORDER BY created LIMIT 10", feed.Id)

	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	return items, err

}
