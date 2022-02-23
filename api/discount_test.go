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

func TestCreateProductDiscountAPI(t *testing.T) {
	admin, _ := randomPDSuperAdmin(t)
	discount := createRandomProductDiscount(t)

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
				"name":             discount.Name,
				"description":      discount.Description,
				"discount_percent": discount.DiscountPercent,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateDiscountParams{
					Name:            discount.Name,
					Description:     discount.Description,
					DiscountPercent: discount.DiscountPercent,
				}

				store.EXPECT().
					CreateDiscount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(discount, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchProductDiscount(t, recorder.Body, discount)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"name":             discount.Name,
				"description":      discount.Description,
				"discount_percent": discount.DiscountPercent,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, false, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateDiscountParams{
					Name:            discount.Name,
					Description:     discount.Description,
					DiscountPercent: discount.DiscountPercent,
				}

				store.EXPECT().
					CreateDiscount(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"name":             discount.Name,
				"description":      discount.Description,
				"discount_percent": discount.DiscountPercent,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateDiscountParams{
					Name:            discount.Name,
					Description:     discount.Description,
					DiscountPercent: discount.DiscountPercent,
				}

				store.EXPECT().
					CreateDiscount(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name":             discount.Name,
				"description":      discount.Description,
				"discount_percent": discount.DiscountPercent,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateDiscountParams{
					Name:            discount.Name,
					Description:     discount.Description,
					DiscountPercent: discount.DiscountPercent,
				}

				store.EXPECT().
					CreateDiscount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Discount{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidDiscount",
			body: gin.H{
				"name":             discount.Name,
				"description":      discount.Description,
				"discount_percent": "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDiscount(gomock.Any(), gomock.Any()).
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

			url := "/discounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetProductDiscountAPI(t *testing.T) {
	productDiscount := createRandomProductDiscount(t)

	testCases := []struct {
		name          string
		DiscountID    int64
		body          gin.H
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id": productDiscount.ID,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDiscount(gomock.Any(), gomock.Eq(productDiscount.ID)).
					Times(1).
					Return(productDiscount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchProductDiscount(t, recorder.Body, productDiscount)
			},
		},
		{
			name:       "NotFound",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id": productDiscount.ID,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDiscount(gomock.Any(), gomock.Eq(productDiscount.ID)).
					Times(1).
					Return(db.Discount{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id": productDiscount.ID,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDiscount(gomock.Any(), gomock.Eq(productDiscount.ID)).
					Times(1).
					Return(db.Discount{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			DiscountID: 0,
			body: gin.H{
				"id": 0,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDiscount(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/discounts/%d", tc.DiscountID)
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkResponse(t, recorder)
		})

	}

}

func TestListProductDiscountAPI(t *testing.T) {
	n := 5
	productdiscounts := make([]db.Discount, n)
	for i := 0; i < n; i++ {
		productdiscounts[i] = createRandomProductDiscount(t)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListDiscountsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListDiscounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(productdiscounts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchProductDiscounts(t, recorder.Body, productdiscounts)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListDiscounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Discount{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListDiscounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListDiscounts(gomock.Any(), gomock.Any()).
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

			url := "/discounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateProductDiscountAPI(t *testing.T) {
	admin, _ := randomPISuperAdmin(t)
	productDiscount := createRandomProductDiscount(t)

	testCases := []struct {
		name          string
		DiscountID    int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id":     productDiscount.ID,
				"active": false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateDiscountParams{
					ID:     productDiscount.ID,
					Active: false,
				}
				store.EXPECT().
					UpdateDiscount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(productDiscount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "Unauthorized",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id":     productDiscount.ID,
				"active": false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, false, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDiscount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "NoAuthorization",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id":     productDiscount.ID,
				"active": false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDiscount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id":     productDiscount.ID,
				"active": false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateDiscountParams{
					ID:     productDiscount.ID,
					Active: false,
				}
				store.EXPECT().
					UpdateDiscount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Discount{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			DiscountID: 0,
			body: gin.H{
				"id":     0,
				"active": false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDiscount(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/discounts/%d", tc.DiscountID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteProductDiscountAPI(t *testing.T) {
	admin, _ := randomPISuperAdmin(t)
	productDiscount := createRandomProductDiscount(t)

	testCases := []struct {
		name          string
		body          gin.H
		DiscountID    int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id": productDiscount.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDiscount(gomock.Any(), gomock.Eq(productDiscount.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "Unauthorized",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id": productDiscount.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, false, time.Minute)

			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDiscount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "No Authorization",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id": productDiscount.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDiscount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "NotFound",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id": productDiscount.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDiscount(gomock.Any(), gomock.Eq(productDiscount.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			DiscountID: productDiscount.ID,
			body: gin.H{
				"id": productDiscount.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDiscount(gomock.Any(), gomock.Eq(productDiscount.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			DiscountID: 0,
			body: gin.H{
				"id": 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationForAdmin(t, request, tokenMaker, authorizationTypeBearer, admin.ID, admin.Username, admin.TypeID, admin.Active, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDiscount(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/discounts/%d", tc.DiscountID)
			request, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkResponse(t, recorder)
		})

	}

}

func randomPDSuperAdmin(t *testing.T) (admin db.Admin, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	admin = db.Admin{
		ID:       util.RandomMoney(),
		Username: util.RandomUser(),
		Email:    util.RandomEmail(),
		Password: hashedPassword,
		Active:   true,
		TypeID:   1,
	}
	return
}

func createRandomProductDiscount(t *testing.T) (discount db.Discount) {
	discount = db.Discount{
		ID:              util.RandomInt(1, 10),
		Name:            util.RandomUser(),
		Description:     util.RandomUser(),
		DiscountPercent: fmt.Sprint(util.RandomMoney()),
		Active:          true,
	}
	return
}

func requireBodyMatchProductDiscount(t *testing.T, body *bytes.Buffer, discount db.Discount) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotDiscount db.Discount
	err = json.Unmarshal(data, &gotDiscount)

	require.NoError(t, err)
	require.Equal(t, discount.ID, gotDiscount.ID)
	require.Equal(t, discount.Name, gotDiscount.Name)
	require.Equal(t, discount.Description, gotDiscount.Description)
	require.Equal(t, discount.DiscountPercent, gotDiscount.DiscountPercent)
	require.Equal(t, discount.Active, gotDiscount.Active)
}

func requireBodyMatchProductDiscounts(t *testing.T, body *bytes.Buffer, discounts []db.Discount) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotDiscounts []db.Discount
	err = json.Unmarshal(data, &gotDiscounts)
	require.NoError(t, err)
	require.Equal(t, discounts, gotDiscounts)
}
