package models

type NotificationType string

const (
	NotificationTypeMention       NotificationType = "mention"
	NotificationTypeStatus        NotificationType = "status"
	NotificationTypeReblog        NotificationType = "reblog"
	NotificationTypeFollow        NotificationType = "follow"
	NotificationTypeFollowRequest NotificationType = "follow_request "
	NotificationTypeFavourite     NotificationType = "favourite"
	NotificationTypePoll          NotificationType = "poll"
	NotificationTypeUpdate        NotificationType = "update"
	NotificationTypeAdminSignUp   NotificationType = "admin.sign_up"
	NotificationTypeAdminReport   NotificationType = "admin.report"

	NotificationTypeUnknown NotificationType = "unknown"
)

type Notification struct {
	Id        string           `json:"id"`
	Type      NotificationType `json:"type"`
	CreatedAt string           `json:"created_at"`
	Account   Account          `json:"account"`
	Status    *Status          `json:"status,omitempty"`
	// FIXME: not implemented
	Report any `json:"report,omitempty"`
}

func (t NotificationType) ToMkNotificationType() MkNotificationType {
	switch t {
	case NotificationTypeStatus:
		return MkNotificationTypeNote
	case NotificationTypeFollow:
		return MkNotificationTypeFollow
	case NotificationTypeFollowRequest:
		return MkNotificationTypeReceiveReaction
	case NotificationTypeFavourite:
		return MkNotificationTypeReceiveReaction
	case NotificationTypeReblog:
		return MkNotificationTypeReceiveRenote
	case NotificationTypeMention:
		return MkNotificationTypeMention
	default:
		return MkNotificationTypeUnknown
	}
}
