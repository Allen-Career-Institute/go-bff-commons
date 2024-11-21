package commons

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetSelectedBatchIds_Empty(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?selected_batch_list=", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	result := GetSelectedBatchIds(ctx)
	assert.Equal(t, []string{}, result)
}

func TestGetSelectedBatchIds_WithBatchIds(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?selected_batch_list=batch1,batch2,batch3", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	result := GetSelectedBatchIds(ctx)
	assert.Equal(t, []string{"batch1", "batch2", "batch3"}, result)
}

func TestGetSelectedBatchIds_WithEmptyStrings(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?selected_batch_list=batch1,,batch3", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	result := GetSelectedBatchIds(ctx)
	assert.Equal(t, []string{"batch1", "batch3"}, result)
}

func TestGetBatchId_Empty(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?batch_id=", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	result := GetBatchID(ctx)

	assert.Equal(t, EmptyStr, result)
}

func TestGetBatchId_WithBatchIds(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?batch_id=batch1,batch2,batch3", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	result := GetBatchID(ctx)

	assert.Equal(t, "batch1", result)
}

func TestGetBatchId_NoBatchIds(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?batch_id=", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	result := GetBatchID(ctx)

	assert.Equal(t, EmptyStr, result)
}
