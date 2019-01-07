package users

import (
	"context"
	"time"
)

type User struct {
	Id        string    `datastore:"-" json:"id" `
	Name      string    `datastore:",noindex" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `datastore:",noindex" json:"updatedAt"`
	// Key *datastore.Key `datastore:"__key__" json:"-"`
}

type IUserRepository interface {
	Create(ctx context.Context, user *User) error

	Find(ctx context.Context, id string) (*User, error)

	List(ctx context.Context) ([]*User, error)

	Delete(ctx context.Context, id string) error

	Update(ctx context.Context, user *User) error
}
