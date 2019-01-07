package users

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"google.golang.org/appengine/aetest"
)

func resetDatastore(ctx context.Context, t *testing.T) {
	repository := UserDatastoreRepository{}
	userList, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("err:%v", err)
	}

	if len(userList) == 0 {
		// tLogWithMessage(t, "Done with reset datastore")
		return
	}

	err = repository.DeleteMulti(ctx, userList)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	resetDatastore(ctx, t)
}
func TestUserDatastoreRepository(t *testing.T) {
	// <setup code>
	log.Print("Setup	AppEngine	Context")
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Errorf("err:%v", err)
	}
	defer done()

	// err = datastore.RunInTransaction(ctx, func(tc context.Context) error {
	// 	return nil
	// }, nil)

	// Create
	testRun(ctx, t, "Create_WhenPassingEmptyId_ReturnError", func(t *testing.T) {
		user := newDummyUserWithEmptyId()
		test_Create_WhenPassingInvalidUser_ReturnErr(ctx, t, user)
	})

	testRun(ctx, t, "Create_WhenPassingEmptyName_ReturnError", func(t *testing.T) {
		user := newDummyUser()
		user.Name = ""
		test_Create_WhenPassingInvalidUser_ReturnErr(ctx, t, user)
	})

	testRun(ctx, t, "Create_WhenPassingValidUser_ReturnNonError", func(t *testing.T) {
		test_Create_WhenPassingValidUser_ReturnNonErr(ctx, t)
	})

	// CreateMulti
	testRun(ctx, t, "CreateMulti_WhenPassingEmptyIds_ReturnError", func(t *testing.T) {
		test_CreateMulti_WhenPassingEmptyIds_ReturnError(ctx, t)
	})

	testRun(ctx, t, "CreateMulti_WhenPassingInvalidIds_ReturnError", func(t *testing.T) {
		test_CreateMulti_WhenPassingInvalidIds_ReturnError(ctx, t)
	})

	testRun(ctx, t, "CreateMulti_WhenPassingValidUserList_ReturnNonError", func(t *testing.T) {
		test_CreateMulti_WhenPassingValidUserList_ReturnNonError(ctx, t)
	})

	// Find
	testRun(ctx, t, "Find_WhenPassingNotExistingId_ReturnError", func(t *testing.T) {
		test_Find_WithNotExistingId_ReturnErr(ctx, t)
	})

	testRun(ctx, t, "Find_WhenPassingExistingId_ReturnUser", func(t *testing.T) {
		test_Find_WhenPassingExistingId_ReturnTheUser(ctx, t)
	})

	// FindMulti
	testRun(ctx, t, "FindMulti_WhenPassingEmptyIds_ReturnError", func(t *testing.T) {
		test_FindMulti_WhenPassingEmptyIds_ReturnError(ctx, t)
	})

	testRun(ctx, t, "FindMulti_WhenPassingInvalidIds_ReturnError", func(t *testing.T) {
		test_FindMulti_WhenPassingInvalidIds_ReturnError(ctx, t)
	})

	testRun(ctx, t, "FindMulti_WhenPassingValidUserList_ReturnTheUserList", func(t *testing.T) {
		test_FindMulti_WhenPassingValidUserList_ReturnTheUserList(ctx, t)
	})

	// List
	testRun(ctx, t, "List", func(t *testing.T) {
		test_List_ReturnUserList(ctx, t)
	})

	// Delete
	testRun(ctx, t, "Delete_WhenPassingNonExistingUser_ReturnError", func(t *testing.T) {
		test_Delete_WhenPassingNonExistingUser_ReturnError(ctx, t)
	})

	testRun(ctx, t, "Delete_WhenPassingExistingUser_ReturnNonError", func(t *testing.T) {
		test_Delete_WhenPassingExistingUser_ReturnNonError(ctx, t)
	})

	// DeleteMulti
	testRun(ctx, t, "DeleteMulti_WhenPassingEmptyIds_ReturnError", func(t *testing.T) {
		test_DeleteMulti_WhenPassingEmptyIds_ReturnError(ctx, t)
	})

	testRun(ctx, t, "DeleteMulti_WhenPassingInvalidIds_ReturnError", func(t *testing.T) {
		test_DeleteMulti_WhenPassingInvalidIds_ReturnError(ctx, t)
	})

	testRun(ctx, t, "DeleteMulti_WhenPassingValidUserList_ReturnNonError", func(t *testing.T) {
		test_DeleteMulti_WhenPassingValidUserList_ReturnNonError(ctx, t)
	})

	// Update
	testRun(ctx, t, "Update_WhenPassingEmptyId_ReturnError", func(t *testing.T) {
		test_Update_WhenPassingEmptyId_ReturnError(ctx, t)
	})

	testRun(ctx, t, "Update_WhenPassingChangedNameUser_ReturnUpdatedUser", func(t *testing.T) {
		test_Update_WhenPassingChangedNameUser_ReturnUpdatedUser(ctx, t)
	})

	// <tear-down code>
}

func testRun(ctx context.Context, t *testing.T, name string, f func(t *testing.T)) {
	// setup
	resetDatastore(ctx, t)
	t.Run(name, f)
}

func test_CreateMulti_WhenPassingEmptyIds_ReturnError(ctx context.Context, t *testing.T) {
	var userList []*User
	repository := UserDatastoreRepository{}
	err := repository.CreateMulti(ctx, userList)
	if err == nil {
		t.Errorf("Error must be thrown")
	}
}

func test_CreateMulti_WhenPassingInvalidIds_ReturnError(ctx context.Context, t *testing.T) {

	var userList []*User
	for i := 0; i < 10; i++ {
		dummyUser := newDummyUserWithEmptyId()
		userList = append(userList, dummyUser)
	}

	repository := UserDatastoreRepository{}
	err := repository.CreateMulti(ctx, userList)
	if err == nil {
		t.Errorf("Error must be thrown")
	}
}

func test_CreateMulti_WhenPassingValidUserList_ReturnNonError(ctx context.Context, t *testing.T) {

	var userList []*User
	for i := 0; i < 10; i++ {
		dummyUser := newDummyUser()
		userList = append(userList, dummyUser)
	}

	repository := UserDatastoreRepository{}
	err := repository.CreateMulti(ctx, userList)
	if err != nil {
		t.Errorf("err:%v", err)
	}

	for _, u := range userList {
		foundUser, err := repository.Find(ctx, u.Id)
		if err != nil {
			t.Errorf("err:%v", err)
		}

		if foundUser == nil {
			t.Errorf("Found user should not be nil")
		}
		// log.Printf("user	%v", foundUser)
	}
}

func test_Create_WhenPassingInvalidUser_ReturnErr(ctx context.Context, t *testing.T, user *User) {
	repository := UserDatastoreRepository{}
	err := repository.Create(ctx, user)
	if err == nil {
		t.Errorf("Error must be thrown")
	}
}

func test_Create_WhenPassingValidUser_ReturnNonErr(ctx context.Context, t *testing.T) {
	repository := UserDatastoreRepository{}
	user := newDummyUser()
	err := repository.Create(ctx, user)
	if err != nil {
		t.Errorf("err:%v", err)
	}
}

func test_Find_WhenPassingExistingId_ReturnTheUser(ctx context.Context, t *testing.T) {
	repository := UserDatastoreRepository{}

	user := newDummyUser()

	createDummyUser(ctx, t, user)

	foundUser, err := repository.Find(ctx, user.Id)

	if err != nil {
		t.Errorf("err:%v", err)
	}

	if foundUser.Id != user.Id {
		t.Errorf("Founded user must be the same with created user	user:%v	foundedUser:%v", user, foundUser)
	}

}

func test_Find_WithNotExistingId_ReturnErr(ctx context.Context, t *testing.T) {
	repository := UserDatastoreRepository{}
	_, err := repository.Find(ctx, "")
	if err == nil {
		t.Errorf("User should be null")
	}
}

func test_FindMulti_WhenPassingEmptyIds_ReturnError(ctx context.Context, t *testing.T) {
	var ids []string
	repository := UserDatastoreRepository{}
	foundUserList, err := repository.FindMulti(ctx, ids)

	if err == nil {
		t.Errorf("Error must be thrown")
	}

	if len(foundUserList) > 0 {
		t.Errorf("FoundUserList should be empty")
	}
}

func test_FindMulti_WhenPassingInvalidIds_ReturnError(ctx context.Context, t *testing.T) {
	var ids []string
	for i := 0; i < 10; i++ {
		ids = append(ids, "")
	}

	repository := UserDatastoreRepository{}
	foundUserList, err := repository.FindMulti(ctx, ids)
	if err == nil {
		t.Errorf("Error must be thrown")
	}

	if len(foundUserList) > 0 {
		t.Errorf("FoundUserList should be empty")
	}
}

func test_FindMulti_WhenPassingValidUserList_ReturnTheUserList(ctx context.Context, t *testing.T) {
	userList := setupDummyUserList(ctx, t)

	var ids []string
	for _, u := range userList {
		ids = append(ids, u.Id)
	}

	repository := UserDatastoreRepository{}
	foundUserList, err := repository.FindMulti(ctx, ids)
	if err != nil {
		t.Errorf("err:%v", err)
	}

	if len(foundUserList) != len(userList) {
		t.Errorf("FoundUserList should have the same length slice	foundUserList:%#v	userList:%#v", foundUserList, userList)
	}
}

func test_Delete_WhenPassingNonExistingUser_ReturnError(ctx context.Context, t *testing.T) {
	dummyId := uuid.New().String()
	repository := UserDatastoreRepository{}

	// user1, err := repository.Find(ctx, dummyId)
	// log.Printf("user:%v", user1)

	err := repository.Delete(ctx, dummyId)
	if err == nil {
		t.Errorf("Error must be thrown")
	}
}

func test_Delete_WhenPassingExistingUser_ReturnNonError(ctx context.Context, t *testing.T) {
	dummyUser := newDummyUser()
	dummyId := dummyUser.Id

	createDummyUser(ctx, t, dummyUser)

	repository := UserDatastoreRepository{}

	// user1, err := repository.Find(ctx, dummyId)
	// log.Printf("user:%v", user1)

	err := repository.Delete(ctx, dummyId)
	if err != nil {
		t.Errorf("err:%v", err)
	}

	user, err := repository.Find(ctx, dummyId)
	if user != nil {
		t.Errorf("User should be not found	user:%v", user)
	}
}

func test_DeleteMulti_WhenPassingValidUserList_ReturnNonError(ctx context.Context, t *testing.T) {

	userList := setupDummyUserList(ctx, t)

	repository := UserDatastoreRepository{}
	err := repository.DeleteMulti(ctx, userList)
	if err != nil {
		t.Errorf("err:%v", err)
	}

	// todo findMulti
}

func test_DeleteMulti_WhenPassingEmptyIds_ReturnError(ctx context.Context, t *testing.T) {

	var userList []*User
	repository := UserDatastoreRepository{}
	err := repository.DeleteMulti(ctx, userList)
	if err == nil {
		t.Errorf("Error must be thrown")
	}
}

func test_DeleteMulti_WhenPassingInvalidIds_ReturnError(ctx context.Context, t *testing.T) {

	var userList []*User
	for i := 0; i < 10; i++ {
		dummyUser := newDummyUserWithEmptyId()
		userList = append(userList, dummyUser)
	}

	repository := UserDatastoreRepository{}
	err := repository.DeleteMulti(ctx, userList)
	if err == nil {
		t.Errorf("Error must be thrown")
	}
}

func test_Update_WhenPassingEmptyId_ReturnError(ctx context.Context, t *testing.T) {
	repository := UserDatastoreRepository{}
	dummyUser := newDummyUserWithEmptyId()
	err := repository.Update(ctx, dummyUser)
	if err == nil {
		t.Errorf("Error must be thrown")
	}
}

func test_Update_WhenPassingChangedNameUser_ReturnUpdatedUser(ctx context.Context, t *testing.T) {
	repository := UserDatastoreRepository{}
	dummyUser := newDummyUser()
	createDummyUser(ctx, t, dummyUser)
	// log.Printf("dummyUser:%#v", dummyUser)

	updatedUser := &User{
		Id:        dummyUser.Id,
		Name:      "ChangedName",
		CreatedAt: dummyUser.CreatedAt,
	}

	err := repository.Update(ctx, updatedUser)
	if err != nil {
		t.Errorf("err:%v", err)
	}
	// log.Printf("updatedUser:%#v", updatedUser)

	foundUser, err := repository.Find(ctx, dummyUser.Id)
	if err != nil {
		t.Errorf("err:%v", err)
	}

	// log.Printf("foundUser:%#v", foundUser)
	if foundUser.Name != updatedUser.Name {
		t.Errorf("User name must be updated updatedUser:%#v	foundUser:%#v", updatedUser, foundUser)
	}
}

func test_List_ReturnUserList(ctx context.Context, t *testing.T) {
	setupDummyUserList(ctx, t)
	repository := UserDatastoreRepository{}
	users, err := repository.List(ctx)
	if err != nil {
		t.Errorf("err:%v", err)
	}
	// tLogWithMessage(t, fmt.Sprintf("users:%#v", users))
	if len(users) == 0 {
		t.Errorf("User List should not be empty")
	}
}

func setupDummyUserList(ctx context.Context, t *testing.T) []*User {
	userList := newDummyUserList()
	createDummyUsers(ctx, t, userList)
	return userList
}

func createDummyUser(ctx context.Context, t *testing.T, u *User) {
	var userList []*User
	userList = append(userList, u)
	createDummyUsers(ctx, t, userList)
}

func createDummyUsers(ctx context.Context, t *testing.T, userList []*User) {
	repository := UserDatastoreRepository{}
	err := repository.CreateMulti(ctx, userList)
	if err != nil {
		t.Errorf("err:%v", err)
	}
}

func newDummyUser() *User {
	return &User{
		Id:        uuid.New().String(),
		Name:      fake.FirstName(),
		CreatedAt: time.Now(),
	}
}

func newDummyUserWithEmptyId() *User {
	u := newDummyUser()
	u.Id = ""
	return u
}

func newDummyUserList() []*User {
	var userList []*User
	for i := 0; i < 10; i++ {
		dummyUser := newDummyUser()
		userList = append(userList, dummyUser)
	}
	return userList
}
