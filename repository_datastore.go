package users

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/appengine/datastore"
)

type UserDatastoreRepository struct {
}

var _ IUserRepository = &UserDatastoreRepository{}

const (
	kind = "User"
)

func newKey(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, kind, id, 0, nil)
}

func newKeys(ctx context.Context, userList []*User) ([]*datastore.Key, error) {
	var keys []*datastore.Key
	for _, u := range userList {
		err := isValidUser(u)
		if err != nil {
			return nil, err
		}
		keys = append(keys, newKey(ctx, u.Id))
	}
	return keys, nil
}

func newKeysByIds(ctx context.Context, ids []string) ([]*datastore.Key, error) {
	var keys []*datastore.Key
	for _, id := range ids {
		if id == "" {
			return nil, fmt.Errorf("datastore: id can not be empty")
		}
		keys = append(keys, newKey(ctx, id))
	}
	return keys, nil
}

func isValidUser(u *User) error {
	if u.Id == "" {
		return fmt.Errorf("datastore: user id empty User: %v", u)
	}

	if u.Name == "" {
		return fmt.Errorf("datastore: user name empty User: %v", u)
	}
	return nil
}

func (repository *UserDatastoreRepository) Create(ctx context.Context, user *User) error {

	err := isValidUser(user)
	if err != nil {
		return err
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	key := datastore.NewKey(ctx, kind, user.Id, 0, nil)
	key, err = datastore.Put(ctx, key, user)
	if err != nil {
		return fmt.Errorf("datastore: could not create User: %v	err:%v", user, err)
	}

	return nil
}

func (repository *UserDatastoreRepository) CreateMulti(ctx context.Context, userList []*User) error {

	if len(userList) == 0 {
		return fmt.Errorf("datastore: userList can not be empty")
	}

	var keys []*datastore.Key
	for _, u := range userList {
		err := isValidUser(u)
		if err != nil {
			return err
		}
		keys = append(keys, newKey(ctx, u.Id))
	}
	// log.Printf("CreateMulti	keys:%v", keys)

	keys, err := datastore.PutMulti(ctx, keys, userList)

	if err != nil {
		return err
	}
	// log.Printf("CreateMulti	userList:%v", userList)

	return nil
}

func (repository *UserDatastoreRepository) Find(ctx context.Context, id string) (*User, error) {
	key := datastore.NewKey(ctx, kind, id, 0, nil)
	user := &User{}
	if err := datastore.Get(ctx, key, user); err != nil {
		return nil, fmt.Errorf("datastore: could not find User	id:%s	err: %v", id, err)
	}
	user.Id = key.StringID()
	return user, nil
}

func (repository *UserDatastoreRepository) FindMulti(ctx context.Context, ids []string) ([]*User, error) {

	if len(ids) == 0 {
		return nil, fmt.Errorf("datastore: ids can not be empty")
	}

	keys, err := newKeysByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	// log.Printf("DeleteMulti	keys:%v", keys)

	var userList = make([]*User, len(keys))

	err = datastore.GetMulti(ctx, keys, userList)
	if err != nil {
		return nil, err
	}

	return userList, nil
}

func (repository *UserDatastoreRepository) Delete(ctx context.Context, id string) error {
	user, err := repository.Find(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("datastore: user doesn't exist id%s", id)
	}

	key := datastore.NewKey(ctx, kind, id, 0, nil)
	err = datastore.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("datastore: could not delete User	id:%s	err: %v", id, err)
	}
	return nil
}

func (repository *UserDatastoreRepository) DeleteMulti(ctx context.Context, userList []*User) error {

	if len(userList) == 0 {
		return fmt.Errorf("datastore: userList can not be empty")
	}

	keys, err := newKeys(ctx, userList)
	if err != nil {
		return err
	}

	// log.Printf("DeleteMulti	keys:%v", keys)

	err = datastore.DeleteMulti(ctx, keys)

	if err != nil {
		return err
	}

	return nil
}

func (repository *UserDatastoreRepository) Update(ctx context.Context, user *User) error {
	if user.Id == "" {
		return fmt.Errorf("user id empty User: %v", user)
	}
	key := datastore.NewKey(ctx, kind, user.Id, 0, nil)
	user.UpdatedAt = time.Now()
	key, err := datastore.Put(ctx, key, user)
	if err != nil {
		return fmt.Errorf("datastore: could not update User: %v	err:%v", user, err)
	}
	return nil
}

func (repository *UserDatastoreRepository) List(ctx context.Context) ([]*User, error) {
	q := datastore.NewQuery(kind).Order("-CreatedAt").Limit(20)
	var users []*User
	keys, err := q.GetAll(ctx, &users)
	if err != nil {
		return nil, fmt.Errorf("datastore: could not retrieve User list	Err:%v", err)
	}

	for i := 0; i < len(keys); i++ {
		// log.Printf("%#v", keys[i])
		users[i].Id = keys[i].StringID()
	}

	// log.Printf("%#v", users)
	return users, nil
}
