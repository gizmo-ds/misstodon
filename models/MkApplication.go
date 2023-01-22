package models

// https://github.com/misskey-dev/misskey/blob/26ae2dfc0f494c377abd878c00044049fcd2bf37/packages/backend/src/misc/api-permissions.ts#L1
const (
	MkAppPermissionReadAccount        = "read:account"
	MkAppPermissionWriteAccount       = "write:account"
	MkAppPermissionReadBlocks         = "read:blocks"
	MkAppPermissionWriteBlocks        = "write:blocks"
	MkAppPermissionReadDrive          = "read:drive"
	MkAppPermissionWriteDrive         = "write:drive"
	MkAppPermissionReadFavorites      = "read:favorites"
	MkAppPermissionWriteFavorites     = "write:favorites"
	MkAppPermissionReadFollowing      = "read:following"
	MkAppPermissionWriteFollowing     = "write:following"
	MkAppPermissionReadMessaging      = "read:messaging"
	MkAppPermissionWriteMessaging     = "write:messaging"
	MkAppPermissionReadMutes          = "read:mutes"
	MkAppPermissionWriteMutes         = "write:mutes"
	MkAppPermissionWriteNotes         = "write:notes"
	MkAppPermissionReadNotifications  = "read:notifications"
	MkAppPermissionWriteNotifications = "write:notifications"
	MkAppPermissionReadReactions      = "read:reactions"
	MkAppPermissionWriteReactions     = "write:reactions"
	MkAppPermissionWriteVotes         = "write:votes"
	MkAppPermissionReadPages          = "read:pages"
	MkAppPermissionWritePages         = "write:pages"
	MkAppPermissionReadPageLikes      = "read:page-likes"
	MkAppPermissionWritePageLikes     = "write:page-likes"
	MkAppPermissionReadUserGroups     = "read:user-groups"
	MkAppPermissionWriteUserGroups    = "write:user-groups"
	MkAppPermissionReadChannels       = "read:channels"
	MkAppPermissionWriteChannels      = "write:channels"
	MkAppPermissionReadGallery        = "read:gallery"
	MkAppPermissionWriteGallery       = "write:gallery"
	MkAppPermissionReadGalleryLikes   = "read:gallery-likes"
	MkAppPermissionWriteGalleryLikes  = "write:gallery-likes"
)

type MkApplication struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	CallbackUrl  string   `json:"callbackUrl"`
	Permission   []string `json:"permission"`
	Secret       string   `json:"secret"`
	IsAuthorized bool     `json:"isAuthorized"`
}
