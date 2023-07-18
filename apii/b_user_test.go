package apii

import (
	"bytes"
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
	return sqlc.User{
		Username:       random.RandomUsername(),
		FullName:       random.RandomFullName(),
		HashedPassword: random.RandomHashedPassword(),
		Email:          random.RandomEmail(),
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStore(ctrl)

	// * build stubs

	// * I expect the GetUser function of the store to be called with any context
	// * and this specific username.Username argument
	store.EXPECT().GetUser(gomock.Any(), gomock.Eq(username.Username)).Times(1).Return(username, nil)

	// * start test server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/users/%s", username.Username)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// * this will send our API request through the server route
	// * and record its response in the recorder
	server.router.ServeHTTP(recorder, request)

	// * check response
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, username)
}
