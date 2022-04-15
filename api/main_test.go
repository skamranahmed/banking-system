package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// run gin in TestMode for tests for cleaner stdout response
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
