package apii

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"simple_shop/db/mock"
	"simple_shop/db/sqlc"
	"simple_shop/db/util/random"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func createRandomUser() sqlc.User {
	username := random.RandomUsername()
	email := username + "@gmail.com"
	return sqlc.User{
		Username:       username,
		FullName:       random.RandomFullName(),
		HashedPassword: random.RandomHashedPassword(),
		Email:          email,
		Level:          false,
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account sqlc.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser sqlc.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, account, gotUser)
}

func TestGetUserAPI(t *testing.T) {
	username := createRandomUser()

	testCase := []struct {
		NameOfTestCase string
		Username       string
		BuildStub      func(store *mock.MockStore)
		CheckResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			NameOfTestCase: "InvalidUsernam",
			Username:       "abc",
			BuildStub: func(store *mock.MockStore) {
				// * build stubs

				// * I expect the GetUser function of the store to be called with any context
				// * and this specific username.Username argument
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// * check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			NameOfTestCase: "NotFound",
			Username:       username.Username,
			BuildStub: func(store *mock.MockStore) {
				// * build stubs

				// * I expect the GetUser function of the store to be called with any context
				// * and this specific username.Username argument
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(username.Username)).Times(1).Return(sqlc.User{}, sql.ErrNoRows)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// * check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			NameOfTestCase: "InternalError",
			Username:       username.Username,
			BuildStub: func(store *mock.MockStore) {
				// * build stubs

				// * I expect the GetUser function of the store to be called with any context
				// * and this specific username.Username argument
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(username.Username)).
					Times(1).
					Return(sqlc.User{}, sql.ErrConnDone)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// * check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			NameOfTestCase: "OK",
			Username:       username.Username,
			BuildStub: func(store *mock.MockStore) {
				// * build stubs

				// * I expect the GetUser function of the store to be called with any context
				// * and this specific username.Username argument
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(username.Username)).
					Times(1).
					Return(username, nil)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// * check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, username)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.NameOfTestCase, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.BuildStub(store)

			// * start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", tc.Username)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// * this will send our API request through the server route
			// * and record its response in the recorder
			server.router.ServeHTTP(recorder, request)

			tc.CheckResponse(t, recorder)
		})
	}
}
