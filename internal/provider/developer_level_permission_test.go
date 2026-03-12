package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandCanViewFinancialDataGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanViewFinancialDataGlobal,
		},
		CanViewFinancialDataGlobal.Expand(),
	)
}

func TestExpandCanManagePermissionsGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanManagePermissionsGlobal,
			CanViewFinancialDataGlobal,
			CanEditGamesGlobal,
			CanPublishGamesGlobal,
			CanReplyToReviewsGlobal,
			CanManagePublicAPKsGlobal,
			CanManageTrackAPKsGlobal,
			CanManageTrackUsersGlobal,
			CanManagePublicListingGlobal,
			CanManageDraftAppsGlobal,
			CanCreateManagedPlayAppsGlobal,
			CanManageOrdersGlobal,
			CanManageAppContentGlobal,
			CanViewNonFinancialDataGlobal,
			CanViewAppQualityGlobal,
			CanManageDeeplinksGlobal,
		},
		CanManagePermissionsGlobal.Expand(),
	)
}

func TestExpandCanEditGamesGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanEditGamesGlobal,
		},
		CanEditGamesGlobal.Expand(),
	)
}

func TestExpandCanPublishGamesGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanPublishGamesGlobal,
		},
		CanPublishGamesGlobal.Expand(),
	)
}

func TestExpandCanReplyToReviewsGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanReplyToReviewsGlobal,
		},
		CanReplyToReviewsGlobal.Expand(),
	)
}

func TestExpandCanManagePublicAPKsGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanManagePublicAPKsGlobal,
		},
		CanManagePublicAPKsGlobal.Expand(),
	)
}

func TestExpandCanManageTrackAPKsGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanManageTrackAPKsGlobal,
		},
		CanManageTrackAPKsGlobal.Expand(),
	)
}

func TestExpandCanManageTrackUsersGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanManageTrackUsersGlobal,
		},
		CanManageTrackUsersGlobal.Expand(),
	)
}

func TestExpandCanManagePublicListingGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanManagePublicListingGlobal,
		},
		CanManagePublicListingGlobal.Expand(),
	)
}

func TestExpandCanManageDraftAppsGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanManageDraftAppsGlobal,
		},
		CanManageDraftAppsGlobal.Expand(),
	)
}

func TestExpandCanCreateManagedPlayAppsGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanCreateManagedPlayAppsGlobal,
		},
		CanCreateManagedPlayAppsGlobal.Expand(),
	)
}

func TestExpandCanChangeManagedPlaySettingsGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{},
		CanChangeManagedPlaySettingGlobal.Expand(),
	)
}

func TestExpandCanManageOrdersGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanManageOrdersGlobal,
		},
		CanManageOrdersGlobal.Expand(),
	)
}

func TestExpandCanManageAppContentGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanManageAppContentGlobal,
		},
		CanManageAppContentGlobal.Expand(),
	)
}

func TestExpandCanViewNonFinancialDataGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanViewNonFinancialDataGlobal,
		},
		CanViewNonFinancialDataGlobal.Expand(),
	)
}

func TestExpandCanViewAppQualityGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanViewAppQualityGlobal,
		},
		CanViewAppQualityGlobal.Expand(),
	)
}

func TestExpandCanManageDeeplinksGlobal(t *testing.T) {
	assert.Equal(
		t,
		[]DeveloperLevelPermission{
			CanManageDeeplinksGlobal,
		},
		CanManageDeeplinksGlobal.Expand(),
	)
}
