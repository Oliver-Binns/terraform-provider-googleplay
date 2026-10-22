package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	androidpublisher "google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

type fakePlayAPI struct {
	existing    []string // permissions of the pre-existing grant; nil means none
	createCalls int
	patchCalls  int
	createCode  int
	patched     []string
}

func (f *fakePlayAPI) handler(t *testing.T) http.Handler {
	const grantName = "developers/dev/users/user@example.com/grants/app1"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/users"):
			user := map[string]any{"email": "user@example.com"}
			if f.existing != nil {
				user["grants"] = []map[string]any{{"name": grantName, "appLevelPermissions": f.existing}}
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"users": []any{user}})
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/grants"):
			f.createCalls++
			if f.createCode != 0 {
				w.WriteHeader(f.createCode)
				_, _ = w.Write([]byte(`{"error":{"code":400,"message":"bad request"}}`))
				return
			}
			var body androidpublisher.Grant
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			body.Name = grantName
			_ = json.NewEncoder(w).Encode(body)
		case r.Method == http.MethodPatch && strings.HasSuffix(r.URL.Path, "/grants/app1"):
			f.patchCalls++
			var body androidpublisher.Grant
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			f.patched = body.AppLevelPermissions
			body.Name = grantName
			_ = json.NewEncoder(w).Encode(body)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func newFakeClient(t *testing.T, f *fakePlayAPI) *GooglePlayClient {
	server := httptest.NewServer(f.handler(t))
	t.Cleanup(server.Close)

	service, err := androidpublisher.NewService(
		context.Background(),
		option.WithEndpoint(server.URL),
		option.WithoutAuthentication(),
	)
	require.NoError(t, err)
	return &GooglePlayClient{service: service, developerID: "dev"}
}

func TestGrantAccessCreatesWhenMissing(t *testing.T) {
	f := &fakePlayAPI{}
	client := newFakeClient(t, f)

	grant, err := client.GrantAccess(context.Background(), "user@example.com", "app1",
		[]AppLevelPermission{"CAN_VIEW_APP_QUALITY"})

	require.NoError(t, err)
	assert.Equal(t, []string{"CAN_VIEW_APP_QUALITY"}, grant.AppLevelPermissions)
	assert.Equal(t, 1, f.createCalls)
	assert.Equal(t, 0, f.patchCalls)
}

func TestGrantAccessReconcilesExistingSamePermissions(t *testing.T) {
	f := &fakePlayAPI{existing: []string{"CAN_VIEW_APP_QUALITY"}}
	client := newFakeClient(t, f)

	grant, err := client.GrantAccess(context.Background(), "user@example.com", "app1",
		[]AppLevelPermission{"CAN_VIEW_APP_QUALITY"})

	require.NoError(t, err)
	assert.Equal(t, []string{"CAN_VIEW_APP_QUALITY"}, grant.AppLevelPermissions)
	assert.Equal(t, 0, f.createCalls)
}

func TestGrantAccessReconcilesExistingDifferentPermissions(t *testing.T) {
	f := &fakePlayAPI{existing: []string{"CAN_VIEW_FINANCIAL_DATA"}}
	client := newFakeClient(t, f)

	grant, err := client.GrantAccess(context.Background(), "user@example.com", "app1",
		[]AppLevelPermission{"CAN_VIEW_APP_QUALITY", "CAN_REPLY_TO_REVIEWS"})

	require.NoError(t, err)
	assert.Equal(t, 0, f.createCalls)
	assert.Equal(t, 1, f.patchCalls)
	assert.ElementsMatch(t, []string{"CAN_VIEW_APP_QUALITY", "CAN_REPLY_TO_REVIEWS"}, f.patched)
	assert.ElementsMatch(t, []string{"CAN_VIEW_APP_QUALITY", "CAN_REPLY_TO_REVIEWS"}, grant.AppLevelPermissions)
}

func TestGrantAccessPropagatesUnrelatedCreateError(t *testing.T) {
	f := &fakePlayAPI{createCode: http.StatusBadRequest}
	client := newFakeClient(t, f)

	_, err := client.GrantAccess(context.Background(), "user@example.com", "app1",
		[]AppLevelPermission{"CAN_VIEW_APP_QUALITY"})

	require.Error(t, err)
	assert.Equal(t, 1, f.createCalls)
	assert.Equal(t, 0, f.patchCalls)
}
