package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"telegram_bot/lib/my_er"
)

type Storage interface {
	Save(page *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(page *Page) error
	IsExists(page *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", my_er.Wrap("cant,t calculate hash", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", my_er.Wrap("cant,t calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
