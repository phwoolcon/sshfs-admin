package session

import (
	"github.com/gin-contrib/sessions"
	gsessions "github.com/gorilla/sessions"
)

type Store interface {
	sessions.Store
}

type store struct {
	*gsessions.FilesystemStore
}

func (c *store) Options(options sessions.Options) {
	c.FilesystemStore.Options = &gsessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

func NewFileStore(path string, keyPairs ...[]byte) Store {
	return &store{gsessions.NewFilesystemStore(path, keyPairs...)}
}
