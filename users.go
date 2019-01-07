package usrsvc

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	Id        string    `datastore:"-" json:"id" `
	Name      string    `datastore:",noindex" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `datastore:",noindex" json:"updatedAt"`
	// Key *datastore.Key `datastore:"__key__" json:"-"`
}

func (u *User) isValid() error {
	if u.Id == "" {
		return fmt.Errorf("datastore: user id empty User: %v", u)
	}

	if u.Name == "" {
		return fmt.Errorf("datastore: user name empty User: %v", u)
	}
	return nil
}

type IUserRepository interface {
	Create(ctx context.Context, user *User) error

	Find(ctx context.Context, id string) (*User, error)

	List(ctx context.Context) ([]*User, error)

	Delete(ctx context.Context, id string) error

	Update(ctx context.Context, user *User) error
}
