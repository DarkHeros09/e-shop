package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/DarkHeros09/e-shop/v2/db/mock"
	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/DarkHeros09/e-shop/v2/token"
	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateShoppingSessionAPI(t *testing.T) {
	user, _ := randomSSUser(t)
	shoppingSession := createRandomShoppingSession(t, user)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_id": shoppingSession.UserID,
				"total":   shoppingSession.Total,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateShoppingSessionParams{
					UserID: shoppingSession.UserID,
					Total:  shoppingSession.Total,
				}

				store.EXPECT().
					CreateShoppingSession(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(shoppingSession, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchShoppingSession(t, recorder.Body, shoppingSession)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"user_id": shoppingSession.UserID,
				"total":   shoppingSession.Total,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateShoppingSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"user_id": shoppingSession.UserID,
				"total":   shoppingSession.Total,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateShoppingSessionParams{
					UserID: shoppingSession.UserID,
					Total:  shoppingSession.Total,
				}

				store.EXPECT().
					CreateShoppingSession(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.ShoppingSession{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "InvalidUserID",
			body: gin.H{
				"user_id": 0,
				"total":   shoppingSession.Total,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateShoppingSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/shopping-sessions"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetShoppingSessionAPI(t *testing.T) {
	user, _ := randomSSUser(t)
	shoppingSession := createRandomShoppingSession(t, user)

	testCases := []struct {
		name              string
		shoppingSessionID int64
		body              gin.H
		setupAuth         func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs        func(store *mockdb.MockStore)
		checkResponse     func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:              "OK",
			shoppingSessionID: shoppingSession.ID,
			body: gin.H{
				"id": shoppingSession.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchShoppingSession(t, recorder.Body, shoppingSession)
			},
		},
		{
			name:              "NoAuthorization",
			shoppingSessionID: shoppingSession.ID,
			body: gin.H{
				"id": shoppingSession.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:              "UnauthorizedUser",
			shoppingSessionID: shoppingSession.ID,
			body: gin.H{
				"id": shoppingSession.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, "unauthorizedUser", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:              "InternalError",
			shoppingSessionID: shoppingSession.ID,
			body: gin.H{
				"id": shoppingSession.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(db.ShoppingSession{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name:              "InvalidShoppingSessionID",
			shoppingSessionID: 0,
			body: gin.H{
				"id": 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/shopping-sessions/%d", tc.shoppingSessionID)
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func randomSSUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:        util.RandomMoney(),
		Username:  util.RandomUser(),
		Password:  hashedPassword,
		Telephone: int32(util.RandomInt(7, 7)),
		Email:     util.RandomEmail(),
	}
	return
}

func createRandomShoppingSession(t *testing.T, user db.User) (shoppingSession db.ShoppingSession) {
	shoppingSession = db.ShoppingSession{
		ID:        util.RandomInt(1, 10),
		UserID:    user.ID,
		Total:     fmt.Sprint(util.RandomMoney()),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	return
}

func requireBodyMatchShoppingSession(t *testing.T, body *bytes.Buffer, shoppingSession db.ShoppingSession) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotShoppingSession db.ShoppingSession
	err = json.Unmarshal(data, &gotShoppingSession)

	require.NoError(t, err)
	require.Equal(t, shoppingSession.UserID, gotShoppingSession.UserID)
	require.Equal(t, shoppingSession.Total, gotShoppingSession.Total)
}
