package base

import (
	"time"
)

type Author struct {
	Code      string    `datastore:"-"`
	Name      string    `datastore:"name,noindex"`
	CreatedAt time.Time `datastore:"created_at"`
	UpdatedAt time.Time `datastore:"updated_at"`
}

type AuthorsBooks struct {
	Code       string    `datastore:"-"`
	Name       string    `datastore:"book,noindex"`
	Href       string    `datastore:"href,noindex"`
	Volume     int       `datastore:"volume"`
	UpdateInfo string    `datastore:"update_info,noindex"`
	CreatedAt  time.Time `datastore:"created_at"`
	UpdatedAt  time.Time `datastore:"updated_at"`
}
