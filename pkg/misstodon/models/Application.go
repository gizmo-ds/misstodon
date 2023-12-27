package models

var (
	ApplicationPermissionRead = []string{
		MkAppPermissionReadAccount,
		MkAppPermissionReadBlocks,
		MkAppPermissionReadFavorites,
		MkAppPermissionReadFollowing,
		MkAppPermissionReadMutes,
		MkAppPermissionReadNotifications,
		MkAppPermissionReadMessaging,
		MkAppPermissionReadDrive,
		MkAppPermissionReadReactions,
	}
	ApplicationPermissionWrite = []string{
		MkAppPermissionWriteAccount,
		MkAppPermissionWriteBlocks,
		MkAppPermissionWriteMessaging,
		MkAppPermissionWriteMutes,
		MkAppPermissionWriteDrive,
		MkAppPermissionWriteNotifications,
		MkAppPermissionWriteNotes,
		MkAppPermissionWriteFavorites,
		MkAppPermissionWriteReactions,
	}
	ApplicationPermissionFollow = []string{
		MkAppPermissionReadBlocks,
		MkAppPermissionWriteBlocks,
		MkAppPermissionWriteFollowing,
		MkAppPermissionReadFollowing,
	}
)

type Application struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Website      *string `json:"website"`
	VapidKey     string  `json:"vapid_key"`
	ClientID     *string `json:"client_id"`
	ClientSecret *string `json:"client_secret"`
	RedirectUri  string  `json:"redirect_uri"`
}
