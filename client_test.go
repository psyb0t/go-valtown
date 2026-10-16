package valtown_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-valtown"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIToken = "test-token"
	testUserID   = "6a1f6e9e-3c4b-4d8a-9b9e-0d6f4c2b1a10"
	testValID    = "0b9c8d7e-6f5a-4b3c-8d2e-1f0a9b8c7d6e"
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *valtown.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := valtown.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testAPIToken)

			return nil
		},
	)

	client, err := valtown.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func assertBearer(t *testing.T, r *http.Request) {
	t.Helper()

	assert.Equal(t, "Bearer "+testAPIToken, r.Header.Get("Authorization"))
}

func TestMeGet(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/me", r.URL.Path)
		assertBearer(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"id":"` + testUserID + `","type":"user","username":"someone",` +
				`"bio":null,"profileImageUrl":null,"tier":"pro",` +
				`"email":"someone@example.com",` +
				`"url":"https://www.val.town/u/someone",` +
				`"links":{"self":"https://api.val.town/v1/me"}}`,
		))
	})

	res, err := client.MeGetWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)

	assert.Equal(t, testUserID, res.JSON200.Id.String())
	assert.Equal(t, "someone", *res.JSON200.Username)
	assert.Equal(
		t,
		valtown.MeGet200JSONResponseBodyTierPro,
		*res.JSON200.Tier,
	)
	assert.Nil(t, res.JSON200.Bio)
}

func TestValsList2(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v2/vals", r.URL.Path)
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "public", r.URL.Query().Get("privacy"))
		assertBearer(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"data":[{"id":"` + testValID + `","name":"hello",` +
				`"privacy":"public","createdAt":"2026-01-02T03:04:05Z"}],` +
				`"links":{"self":"https://api.val.town/v2/vals"}}`,
		))
	})

	privacy := valtown.ValsList2ParamsPrivacyPublic

	res, err := client.ValsList2WithResponse(
		context.Background(),
		&valtown.ValsList2Params{Limit: 10, Privacy: &privacy},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	require.Len(t, res.JSON200.Data, 1)

	val := res.JSON200.Data[0]
	assert.Equal(t, testValID, val.Id.String())
	assert.Equal(t, "hello", val.Name)
	assert.Equal(t, valtown.ValPrivacyPublic, val.Privacy)
}

func TestSqliteExecute(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/sqlite/execute", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assertBearer(t, r)

		var body map[string]any

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "SELECT 1 AS one", body["statement"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"columns":["one"],"columnTypes":[""],` +
				`"rows":[[1]],"rowsAffected":0}`,
		))
	})

	var statement valtown.SqliteExecuteJSONBody_Statement

	require.NoError(t, statement.FromSqliteExecuteJSONBodyStatement0(
		"SELECT 1 AS one",
	))

	res, err := client.SqliteExecuteWithResponse(
		context.Background(),
		valtown.SqliteExecuteJSONRequestBody{Statement: statement},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)

	assert.Equal(t, []string{"one"}, res.JSON200.Columns)
	assert.Equal(t, [][]any{{float64(1)}}, res.JSON200.Rows)
}
