package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/banking-system/config"
	db "github.com/skamranahmed/banking-system/db/sqlc"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	// load config
	config.Load("../config")

	server, err := NewServer(store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	// run gin in TestMode for tests for cleaner stdout response
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
