package provider

type DeveloperLevelPermission string

const (
	UnspecifiedDeveloperLevelPermission DeveloperLevelPermission = "DEVELOPER_LEVEL_PERMISSION_UNSPECIFIED"
	CanViewFinancialDataGlobal          DeveloperLevelPermission = "CAN_VIEW_FINANCIAL_DATA_GLOBAL"
	CanManagePermissionsGlobal          DeveloperLevelPermission = "CAN_MANAGE_PERMISSIONS_GLOBAL"
	CanEditGamesGlobal                  DeveloperLevelPermission = "CAN_EDIT_GAMES_GLOBAL"
	CanPublishGamesGlobal               DeveloperLevelPermission = "CAN_PUBLISH_GAMES_GLOBAL"
	CanReplyToReviewsGlobal             DeveloperLevelPermission = "CAN_REPLY_TO_REVIEWS_GLOBAL"
	CanManagePublicAPKsGlobal           DeveloperLevelPermission = "CAN_MANAGE_PUBLIC_APKS_GLOBAL"
	CanManageTrackAPKsGlobal            DeveloperLevelPermission = "CAN_MANAGE_TRACK_APKS_GLOBAL"
	CanManageTrackUsersGlobal           DeveloperLevelPermission = "CAN_MANAGE_TRACK_USERS_GLOBAL"
	CanManagePublicListingGlobal        DeveloperLevelPermission = "CAN_MANAGE_PUBLIC_LISTING_GLOBAL"
	CanManageDraftAppsGlobal            DeveloperLevelPermission = "CAN_MANAGE_DRAFT_APPS_GLOBAL"
	CanCreateManagedPlayAppsGlobal      DeveloperLevelPermission = "CAN_CREATE_MANAGED_PLAY_APPS_GLOBAL"
	CanChangeManagedPlaySettingGlobal   DeveloperLevelPermission = "CAN_CHANGE_MANAGED_PLAY_SETTING_GLOBAL"
	CanManageOrdersGlobal               DeveloperLevelPermission = "CAN_MANAGE_ORDERS_GLOBAL"
	CanManageAppContentGlobal           DeveloperLevelPermission = "CAN_MANAGE_APP_CONTENT_GLOBAL"
	CanViewNonFinancialDataGlobal       DeveloperLevelPermission = "CAN_VIEW_NON_FINANCIAL_DATA_GLOBAL"
	CanViewAppQualityGlobal             DeveloperLevelPermission = "CAN_VIEW_APP_QUALITY_GLOBAL"
	CanManageDeeplinksGlobal            DeveloperLevelPermission = "CAN_MANAGE_DEEPLINKS_GLOBAL"
)

func (permission DeveloperLevelPermission) Expand() []DeveloperLevelPermission {
	switch permission {
	case CanManagePermissionsGlobal:
		return []DeveloperLevelPermission{
			permission,
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
		}
	case CanChangeManagedPlaySettingGlobal:
		return []DeveloperLevelPermission{}
	default:
		return []DeveloperLevelPermission{permission}
	}
}
