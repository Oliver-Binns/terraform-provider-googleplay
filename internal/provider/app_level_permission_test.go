package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandCanManagePermissions(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanViewFinancialData,
			CanManagePermissions,
			CanReplyToReviews,
			CanManagePublicAPKs,
			CanManageTrackAPKs,
			CanManageTrackUsers,
			CanManagePublicListing,
			CanManageDraftApps,
			CanManageOrders,
			CanManageAppContent,
			CanViewNonFinancialData,
			CanViewAppQuality,
			CanManageDeeplinks,
		},
		CanManagePermissions.Expand(),
	)
}

func TestExpandCanViewFinancialData(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanViewFinancialData,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanViewFinancialData.Expand(),
	)
}

func TestExpandCanReplyToReviews(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanReplyToReviews,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanReplyToReviews.Expand(),
	)
}

func TestExpandCanManagePublicAPKs(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanManagePublicAPKs,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanManagePublicAPKs.Expand(),
	)
}

func TestExpandCanManageTrackAPKs(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanManageTrackAPKs,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanManageTrackAPKs.Expand(),
	)
}

func TestExpandCanManageTrackUsers(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanManageTrackUsers,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanManageTrackUsers.Expand(),
	)
}

func TestExpandCanManagePublicListing(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanManagePublicListing,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanManagePublicListing.Expand(),
	)
}

func TestExpandCanManageDraftApps(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanManageDraftApps,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanManageDraftApps.Expand(),
	)
}

func TestExpandCanManageAppContent(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanManageAppContent,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanManageAppContent.Expand(),
	)
}

func TestExpandCanManageOrders(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanManageOrders,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanManageOrders.Expand(),
	)
}

func TestExpandCanManageDeeplinks(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanManageDeeplinks,
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanManageDeeplinks.Expand(),
	)
}

func TestExpandCanViewNonFinancialData(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanViewNonFinancialData,
			CanViewAppQuality,
		},
		CanViewNonFinancialData.Expand(),
	)
}

func TestExpandCanViewAppQuality(t *testing.T) {
	assert.Equal(
		t,
		[]AppLevelPermission{
			CanViewAppQuality,
		},
		CanViewAppQuality.Expand(),
	)
}
