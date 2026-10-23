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
		case r.Method == http.MethodPatch && strings.Contains(r.URL.Path, "/grants/"):
			f.patchCalls++
			if f.existing == nil {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error":{"code":404,"message":"not found"}}`))
				return
			}
			var body androidpublisher.Grant
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			f.patched = body.AppLevelPermissions
			body.Name = grantName
			body.PackageName = "app1"
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
	assert.Contains(t, err.Error(), "bad request")
	assert.Equal(t, 1, f.createCalls)
	// A reconcile is attempted, but fails as there is no grant, so the create error wins.
	assert.Equal(t, 1, f.patchCalls)
}

func TestGrantMatchesApp(t *testing.T) {
	published := &androidpublisher.Grant{
		Name:        "developers/dev/users/user@example.com/grants/uk.gov.app",
		PackageName: "uk.gov.app",
	}
	draft := &androidpublisher.Grant{
		Name: "developers/dev/users/user@example.com/grants/4973279986054171407",
	}
	malformed := &androidpublisher.Grant{Name: "developers/dev"}

	assert.True(t, grantMatchesApp(published, "uk.gov.app"))
	assert.False(t, grantMatchesApp(published, "4973279986054171407"))
	assert.True(t, grantMatchesApp(draft, "4973279986054171407"))
	assert.False(t, grantMatchesApp(draft, "uk.gov.app"))
	assert.False(t, grantMatchesApp(malformed, "dev"))
}

func TestGrantAccessResolvesNumericIDOfPublishedApp(t *testing.T) {
	// The grant is listed under its package name ("app1"), so the numeric ID
	// does not match the listing and the create is rejected as a duplicate.
	f := &fakePlayAPI{existing: []string{"CAN_VIEW_FINANCIAL_DATA"}, createCode: http.StatusConflict}
	client := newFakeClient(t, f)

	grant, err := client.GrantAccess(context.Background(), "user@example.com", "123456",
		[]AppLevelPermission{"CAN_VIEW_APP_QUALITY"})

	require.NoError(t, err)
	assert.Equal(t, 1, f.createCalls)
	assert.Equal(t, 1, f.patchCalls)
	assert.Equal(t, []string{"CAN_VIEW_APP_QUALITY"}, f.patched)
	assert.Equal(t, "app1", grant.PackageName)
}
