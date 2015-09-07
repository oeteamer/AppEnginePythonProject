package base

import (
	"appengine/datastore"
	"time"
)

type Author struct {
	ID        datastore.Key   `datastore:"-"`
	Code      string          `datastore:"-"`
	Name      string          `datastore:"name,noindex"`
	CreatedAt time.Time       `datastore:"created_at"`
	UpdatedAt time.Time       `datastore:"updated_at"`
	Books     map[string]Book `datastore:"-"`
}

type Book struct {
	ID         datastore.Key `datastore:"-"`
	Code       string        `datastore:"-"`
	Name       string        `datastore:"book,noindex"`
	Href       string        `datastore:"href,noindex"`
	Volume     string        `datastore:"volume"`
	UpdateInfo string        `datastore:"update_info,noindex"`
	CreatedAt  time.Time     `datastore:"created_at"`
	UpdatedAt  time.Time     `datastore:"updated_at"`
}
