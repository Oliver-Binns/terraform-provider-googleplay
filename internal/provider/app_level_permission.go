package provider

type AppLevelPermission string

const (
	CanViewFinancialData    AppLevelPermission = "CAN_VIEW_FINANCIAL_DATA"
	CanManagePermissions    AppLevelPermission = "CAN_MANAGE_PERMISSIONS"
	CanReplyToReviews       AppLevelPermission = "CAN_REPLY_TO_REVIEWS"
	CanManagePublicAPKs     AppLevelPermission = "CAN_MANAGE_PUBLIC_APKS"
	CanManageTrackAPKs      AppLevelPermission = "CAN_MANAGE_TRACK_APKS"
	CanManageTrackUsers     AppLevelPermission = "CAN_MANAGE_TRACK_USERS"
	CanManagePublicListing  AppLevelPermission = "CAN_MANAGE_PUBLIC_LISTING"
	CanManageDraftApps      AppLevelPermission = "CAN_MANAGE_DRAFT_APPS"
	CanManageOrders         AppLevelPermission = "CAN_MANAGE_ORDERS"
	CanManageAppContent     AppLevelPermission = "CAN_MANAGE_APP_CONTENT"
	CanViewNonFinancialData AppLevelPermission = "CAN_VIEW_NON_FINANCIAL_DATA"
	CanViewAppQuality       AppLevelPermission = "CAN_VIEW_APP_QUALITY"
	CanManageDeeplinks      AppLevelPermission = "CAN_MANAGE_DEEPLINKS"
)

func (permission AppLevelPermission) Expand() []AppLevelPermission {
	switch permission {
	case CanManagePermissions:
		return []AppLevelPermission{
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
		}
	case CanViewFinancialData,
		CanReplyToReviews,
		CanManagePublicAPKs,
		CanManageTrackAPKs,
		CanManageTrackUsers,
		CanManagePublicListing,
		CanManageDraftApps,
		CanManageAppContent,
		CanManageOrders,
		CanManageDeeplinks:
		return []AppLevelPermission{
			permission,
			CanViewNonFinancialData,
			CanViewAppQuality,
		}
	case CanViewNonFinancialData:
		return []AppLevelPermission{
			CanViewNonFinancialData,
			CanViewAppQuality,
		}
	default:
		return []AppLevelPermission{permission}
	}
}
