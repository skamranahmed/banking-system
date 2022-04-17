package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/skamranahmed/banking-system/db/mock"
	db "github.com/skamranahmed/banking-system/db/sqlc"
	"github.com/skamranahmed/banking-system/token"
	"github.com/skamranahmed/banking-system/utils"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {
	amount := int64(100)

	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	user3, _ := randomUser(t)

	account1 := randomAccount(uint(user1.ID))
	account2 := randomAccount(uint(user2.ID))
	account3 := randomAccount(uint(user3.ID))

	account1.Currency = utils.INR
	account2.Currency = utils.INR
	account3.Currency = utils.USD

	testCases := []struct {
		name         string
		body         gin.H
		setupAuth    func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs   func(store *mockdb.MockStore)
		expectStatus int
	}{
		{
			name: "Happy Case - All OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.INR,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, uint(user1.ID), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := db.TransferTxnParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTxn(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "Failure Case - FromAccount Does Not Exists / NotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, uint(user1.ID), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().TransferTxn(gomock.Any(), gomock.Any()).Times(0)
			},
			expectStatus: http.StatusNotFound,
		},
		{
			name: "Failure Case - ToAccount Does Not Exists / NotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.INR,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, uint(user1.ID), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().TransferTxn(gomock.Any(), gomock.Any()).Times(0)
			},
			expectStatus: http.StatusNotFound,
		},
		{
			name: "Failure Case - FromAccount Currency Mismatch",
			body: gin.H{
				"from_account_id": account3.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.INR,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, uint(user3.ID), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().TransferTxn(gomock.Any(), gomock.Any()).Times(0)
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Failure Case - ToAccount Currency Mismatch",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        utils.INR,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, uint(user1.ID), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().TransferTxn(gomock.Any(), gomock.Any()).Times(0)
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Failure Case - Unsupported Currency Provided by User",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "XYZ",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, uint(user1.ID), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTxn(gomock.Any(), gomock.Any()).Times(0)
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Failure Case - Negative Amount in Transfer Request",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          -amount,
				"currency":        utils.INR,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, uint(user1.ID), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTxn(gomock.Any(), gomock.Any()).Times(0)
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Failure Case - GetAccountError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.INR,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, uint(user1.ID), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().TransferTxn(gomock.Any(), gomock.Any()).Times(0)
			},
			expectStatus: http.StatusInternalServerError,
		},
		{
			name: "Failure Case - TransferTxnError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.INR,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, uint(user1.ID), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTxn(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxnResult{}, sql.ErrTxDone)
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// server := NewServer(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			// add authorization header to the request
			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, tc.expectStatus, recorder.Code)
		})
	}
}
