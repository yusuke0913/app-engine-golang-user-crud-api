package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"

	"github.com/icrowley/fake"
)

type responseHandlerFunc func(t *testing.T, rr *httptest.ResponseRecorder, apiTest apiTest)

type httpHandlerFunc func(w http.ResponseWriter, r *http.Request)

type setupFunc func(ctx context.Context, t *testing.T, apiTest apiTest)

type apiTest struct {
	name                string
	method              string
	url                 string
	urlVars             map[string]string
	request             requester
	setupFunc           setupFunc
	expectedStatusCode  int
	httpHandlerFunc     httpHandlerFunc
	responseHandlerFunc responseHandlerFunc
}

var apiTests = []apiTest{

	// Create
	{
		name:                "Create_WhenPassingEmptyName_ReturnError",
		method:              "POST",
		url:                 "/users/v1",
		urlVars:             nil,
		request:             userCreateRequest{User: &User{Name: ""}},
		expectedStatusCode:  http.StatusInternalServerError,
		httpHandlerFunc:     createUser,
		responseHandlerFunc: nil,
	},

	{
		name:                "Create_ByUser_ReturnCreatedUser",
		method:              "POST",
		url:                 "/users/v1",
		urlVars:             nil,
		request:             userCreateRequest{User: &User{Name: fake.FirstName()}},
		expectedStatusCode:  http.StatusOK,
		httpHandlerFunc:     createUser,
		responseHandlerFunc: testUserCreateResponse,
	},

	// Find
	{
		name:   "Find_ByNotExistingUser_ReturnError",
		method: "GET",
		url:    "/users/v1/DummyId",
		urlVars: map[string]string{
			"id": "DummyId",
		},
		request:            nil,
		expectedStatusCode: http.StatusInternalServerError,
		httpHandlerFunc:    findUser,
	},

	{
		name:   "Find_ByExistingUser_ReturnTheUser",
		method: "GET",
		url:    "/users/v1/DummyId",
		urlVars: map[string]string{
			"id": "DummyId",
		},
		request:             nil,
		setupFunc:           setupDummyUser,
		expectedStatusCode:  http.StatusOK,
		httpHandlerFunc:     findUser,
		responseHandlerFunc: testUserFindResponse,
	},

	// Update
	{
		name:   "Update_WhenPasingNonExistingUser_ReturnError",
		method: "PUT",
		url:    "/users/v1/DummyId",
		urlVars: map[string]string{
			"id": "DummyId",
		},
		request:             userUpdateRequest{User: &User{Id: "DummyId", Name: "ChangedName"}},
		expectedStatusCode:  http.StatusInternalServerError,
		httpHandlerFunc:     updateUser,
		responseHandlerFunc: nil,
	},

	{
		name:   "Update_WhenPasingExistingUser_ReturnNonError",
		method: "PUT",
		url:    "/users/v1/DummyId",
		urlVars: map[string]string{
			"id": "DummyId",
		},
		request:             userUpdateRequest{User: &User{Id: "DummyId", Name: "ChangedName"}},
		setupFunc:           setupDummyUser,
		expectedStatusCode:  http.StatusOK,
		httpHandlerFunc:     updateUser,
		responseHandlerFunc: testUserUpdateResponse,
	},

	// Delete
	{
		name:   "Delete_WhenPasingNotExistingUser_ReturnError",
		method: "DELETE",
		url:    "/users/v1/DummyId",
		urlVars: map[string]string{
			"id": "DummyId",
		},
		request:            nil,
		expectedStatusCode: http.StatusInternalServerError,
		httpHandlerFunc:    deleteUser,
	},

	{
		name:   "Delete_WhenPasingExistingUser_ReturnNonError",
		method: "DELETE",
		url:    "/users/v1/DummyId",
		urlVars: map[string]string{
			"id": "DummyId",
		},
		request:            nil,
		setupFunc:          setupDummyUser,
		expectedStatusCode: http.StatusOK,
		httpHandlerFunc:    deleteUser,
	},

	// List
	{
		name:                "List_ReturnUserList",
		method:              "GET",
		url:                 "/users/v1/list",
		urlVars:             nil,
		request:             nil,
		setupFunc:           setupDummyUserListWithApiTestCase,
		expectedStatusCode:  http.StatusOK,
		httpHandlerFunc:     listUsers,
		responseHandlerFunc: testUserListResponse,
	},
}

func setupDummyUser(ctx context.Context, t *testing.T, testCase apiTest) {
	user := newDummyUser()
	user.Id = testCase.urlVars["id"]
	createDummyUser(ctx, t, user)
}

func setupDummyUserListWithApiTestCase(ctx context.Context, t *testing.T, testCase apiTest) {
	setupDummyUserList(ctx, t)
}

func TestUsersApiHandler(t *testing.T) {

	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	defer inst.Close()

	req, err := inst.NewRequest("GET", "/users/v1", nil)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	ctx := appengine.NewContext(req)

	for _, tt := range apiTests {
		t.Run(tt.name, func(t *testing.T) {
			resetDatastore(ctx, t)
			if tt.setupFunc != nil {
				tt.setupFunc(ctx, t, tt)
			}
			testApi(t, inst, tt)
		})
	}
}

func testUserCreateResponse(t *testing.T, rr *httptest.ResponseRecorder, apiTest apiTest) {
	req := (apiTest.request).(userCreateRequest)
	var response userCreateResponse
	decodeResponseBody(rr.Body.Bytes(), &response)
	if response.User.Id == "" {
		t.Errorf("User should have ID	response:%v", response)
	}

	if response.User.Name != req.User.Name {
		t.Errorf("User should have the same name	response:%v", response)
	}
}

func testUserFindResponse(t *testing.T, rr *httptest.ResponseRecorder, apiTest apiTest) {
	var response userFindResponse
	decodeResponseBody(rr.Body.Bytes(), &response)

	expectedId := apiTest.urlVars["id"]
	if response.User.Id != expectedId {
		t.Errorf("FoundUserId should be the same with expectedId	expectedId:%v	foundUser:%v", expectedId, response.User)
	}
}

func testUserUpdateResponse(t *testing.T, rr *httptest.ResponseRecorder, apiTest apiTest) {
	req := (apiTest.request).(userUpdateRequest)
	var response userUpdateResponse
	decodeResponseBody(rr.Body.Bytes(), &response)

	if response.User.Name != req.User.Name {
		t.Errorf("User should have the name	response:%v", response)
	}
}

func testUserListResponse(t *testing.T, rr *httptest.ResponseRecorder, apiTest apiTest) {
	var response userListResponse
	decodeResponseBody(rr.Body.Bytes(), &response)

	if len(response.Users) == 0 {
		t.Errorf("UserList should have at least one user	response:%v", response)
	}
}

func testApi(t *testing.T, inst aetest.Instance, apiTest apiTest) {
	body := encodeRequestBody(apiTest.request)

	t.Logf("inst.NewRequest	method:%v	url:%v	urlVars:%v	body:%v", apiTest.method, apiTest.url, apiTest.urlVars, body)
	req, err := inst.NewRequest(apiTest.method, apiTest.url, body)
	if err != nil {
		t.Errorf("err:%v", err)
	}
	req = mux.SetURLVars(req, apiTest.urlVars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiTest.httpHandlerFunc)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != apiTest.expectedStatusCode {
		t.Errorf("handler returned wrong status code: got %v want %v", status, apiTest.expectedStatusCode)
	}

	if apiTest.responseHandlerFunc != nil {
		apiTest.responseHandlerFunc(t, rr, apiTest)
	}
}
