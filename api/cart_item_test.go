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

func TestCreateCartItemAPI(t *testing.T) {
	user, _ := randomCIUser(t)
	shoppingSession := createRandomShoppingSession(t, user)
	cartItem := createRandomCartItem(t, shoppingSession)

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
				"session_id": cartItem.SessionID,
				"product_id": cartItem.ProductID,
				"quantity":   cartItem.Quantity,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				arg := db.CreateCartItemParams{
					SessionID: shoppingSession.ID,
					ProductID: cartItem.ProductID,
					Quantity:  cartItem.Quantity,
				}

				store.EXPECT().
					CreateCartItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(cartItem, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCartItem(t, recorder.Body, cartItem)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"session_id": cartItem.SessionID,
				"product_id": cartItem.ProductID,
				"quantity":   cartItem.Quantity,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreateCartItem(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"session_id": cartItem.SessionID,
				"product_id": cartItem.ProductID,
				"quantity":   cartItem.Quantity,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				store.EXPECT().
					CreateCartItem(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"session_id": cartItem.SessionID,
				"product_id": cartItem.ProductID,
				"quantity":   cartItem.Quantity,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				arg := db.CreateCartItemParams{
					SessionID: cartItem.SessionID,
					ProductID: cartItem.ProductID,
					Quantity:  cartItem.Quantity,
				}

				store.EXPECT().
					CreateCartItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.CartItem{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidSessionID",
			body: gin.H{
				"session_id": 0,
				"product_id": cartItem.ProductID,
				"quantity":   cartItem.Quantity,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateCartItem(gomock.Any(), gomock.Any()).
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

			url := "/cart-items"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetCartItemByIDAPI(t *testing.T) {
	user, _ := randomCIUser(t)
	shoppingSession := createRandomShoppingSessionForCartItem(t, user)
	cartItem := createRandomCartItem(t, shoppingSession)

	testCases := []struct {
		name          string
		SessionID     int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			SessionID: cartItem.SessionID,
			body: gin.H{
				"session_id": shoppingSession.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				store.EXPECT().
					GetCartItemBySessionID(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(cartItem, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCartItem(t, recorder.Body, cartItem)
			},
		},
		{
			name:      "NoAuthorization",
			SessionID: shoppingSession.ID,
			body: gin.H{
				"session_id": shoppingSession.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(0)

				store.EXPECT().
					GetCartItemBySessionID(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "UnauthorizedUser",
			SessionID: shoppingSession.ID,
			body: gin.H{
				"session_id": shoppingSession.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, "unauthorizedUser", time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				store.EXPECT().
					GetCartItemBySessionID(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			SessionID: shoppingSession.ID,
			body: gin.H{
				"session_id": shoppingSession.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				store.EXPECT().
					GetCartItemBySessionID(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(db.CartItem{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			SessionID: shoppingSession.ID,
			body: gin.H{
				"session_id": shoppingSession.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				store.EXPECT().
					GetCartItemBySessionID(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(db.CartItem{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidSessionID",
			SessionID: 0,
			body: gin.H{
				"session_id": 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetCartItemBySessionID(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) // no need to call defer ctrl.finish() in 1.6V

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStub(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/cart-items/%d", tc.SessionID)
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkResponse(t, recorder)
		})

	}

}

// func TestListcart-itemsAPI(t *testing.T) {
// 	n := 5
// 	cart-items := make([]db.CartItem, n)
// 	user, _ := randomCIUser(t)
// 	cartItem1 := createRandomCartItem(t)
// 	cartItem2 := createRandomCartItem(t)
// 	cartItem3 := createRandomCartItem(t)

// 	cart-items = append(cart-items, cartItem1, cartItem2, cartItem3)

// 	type Query struct {
// 		pageID   int
// 		pageSize int
// 	}

// 	testCases := []struct {
// 		name          string
// 		query         Query
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.ListCartItemParams{
// 					Limit:  int32(n),
// 					Offset: 0,
// 				}

// 				store.EXPECT().
// 					ListCartItem(gomock.Any(), gomock.Eq(arg)).
// 					Times(1).
// 					Return(cart-items, nil)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyMatchcart-items(t, recorder.Body, cart-items)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListCartItem(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return([]db.CartItem{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPageID",
// 			query: Query{
// 				pageID:   -1,
// 				pageSize: n,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListCartItem(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPageSize",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: 100000,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListCartItem(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			url := "/cart-items"
// 			request, err := http.NewRequest(http.MethodGet, url, nil)
// 			require.NoError(t, err)

// 			// Add query parameters to request URL
// 			q := request.URL.Query()
// 			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
// 			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
// 			request.URL.RawQuery = q.Encode()

// 			tc.setupAuth(t, request, server.tokenMaker)
// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

func TestUpdateCartItemByUserIDAPI(t *testing.T) {
	user, _ := randomCIUser(t)
	shoppingSession := createRandomShoppingSessionForCartItem(t, user)
	cartItem := createRandomCartItem(t, shoppingSession)

	testCases := []struct {
		name          string
		SessionID     int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			SessionID: shoppingSession.ID,
			body: gin.H{
				"id":         cartItem.ID,
				"session_id": cartItem.SessionID,
				"quantity":   4,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				arg := db.UpdateCartItemParams{
					ID:       cartItem.ID,
					Quantity: 4,
				}

				store.EXPECT().
					UpdateCartItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(cartItem, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			SessionID: shoppingSession.ID,
			body: gin.H{
				"id":         cartItem.ID,
				"session_id": cartItem.SessionID,
				"quantity":   4,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCartItem(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			SessionID: shoppingSession.ID,
			body: gin.H{
				"id":         cartItem.ID,
				"session_id": cartItem.SessionID,
				"quantity":   4,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				arg := db.UpdateCartItemParams{
					ID:       cartItem.ID,
					Quantity: 4,
				}
				store.EXPECT().
					UpdateCartItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.CartItem{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidSessionID",
			SessionID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdateCartItem(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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

			url := fmt.Sprintf("/cart-items/%d", tc.SessionID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteCartItemAPI(t *testing.T) {
	user, _ := randomCIUser(t)
	shoppingSession := createRandomShoppingSessionForCartItem(t, user)
	cartItem := createRandomCartItem(t, shoppingSession)

	testCases := []struct {
		name          string
		ID            int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   cartItem.ID,
			body: gin.H{
				"id":         cartItem.ID,
				"session_id": cartItem.SessionID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				store.EXPECT().
					DeleteCartItem(gomock.Any(), gomock.Eq(cartItem.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   cartItem.ID,
			body: gin.H{
				"id":         cartItem.ID,
				"session_id": cartItem.SessionID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				store.EXPECT().
					DeleteCartItem(gomock.Any(), gomock.Eq(cartItem.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   cartItem.ID,
			body: gin.H{
				"id":         cartItem.ID,
				"session_id": cartItem.SessionID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Eq(shoppingSession.ID)).
					Times(1).
					Return(shoppingSession, nil)

				store.EXPECT().
					DeleteCartItem(gomock.Any(), gomock.Eq(cartItem.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCartItemID",
			ID:   0,
			body: gin.H{
				"id":         0,
				"session_id": cartItem.SessionID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetShoppingSession(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteCartItem(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) // no need to call defer ctrl.finish() in 1.6V

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStub(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/cart-items/%d", tc.ID)
			request, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkResponse(t, recorder)
		})

	}

}

func randomCIUser(t *testing.T) (user db.User, password string) {
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

func createRandomShoppingSessionForCartItem(t *testing.T, user db.User) (shoppingSession db.ShoppingSession) {
	shoppingSession = db.ShoppingSession{
		ID:     util.RandomInt(1, 10),
		UserID: user.ID,
		Total:  fmt.Sprint(util.RandomMoney()),
	}
	return
}

func createRandomCartItem(t *testing.T, shoppingSession db.ShoppingSession) (cartItem db.CartItem) {
	cartItem = db.CartItem{
		ID:        util.RandomMoney(),
		SessionID: shoppingSession.ID,
		ProductID: util.RandomMoney(),
		Quantity:  int32(util.RandomMoney()),
	}
	return
}

func requireBodyMatchCartItem(t *testing.T, body *bytes.Buffer, cartItem db.CartItem) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotCartItem db.CartItem
	err = json.Unmarshal(data, &gotCartItem)

	require.NoError(t, err)
	require.Equal(t, cartItem.ID, gotCartItem.ID)
	require.Equal(t, cartItem.ProductID, gotCartItem.ProductID)
	require.Equal(t, cartItem.SessionID, gotCartItem.SessionID)
	require.Equal(t, cartItem.Quantity, gotCartItem.Quantity)
}

// func requireBodyMatchcart-items(t *testing.T, body *bytes.Buffer, cart-items []db.CartItem) {
// 	data, err := ioutil.ReadAll(body)
// 	require.NoError(t, err)

// 	var gotcart-items []db.CartItem
// 	err = json.Unmarshal(data, &gotcart-items)
// 	require.NoError(t, err)
// 	require.Equal(t, cart-items, gotcart-items)
// }
