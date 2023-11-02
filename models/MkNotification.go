package models

type MkNotificationType string

const (
	MkNotificationTypeNote                  MkNotificationType = "note"
	MkNotificationTypeFollow                MkNotificationType = "follow"
	MkNotificationTypeAchievementEarned     MkNotificationType = "achievementEarned"
	MkNotificationTypeReceiveFollowRequest  MkNotificationType = "receiveFollowRequest"
	MkNotificationTypeFollowRequestAccepted MkNotificationType = "followRequestAccepted"
	MkNotificationTypeReceiveReaction       MkNotificationType = "reaction"
	MkNotificationTypeReceiveRenote         MkNotificationType = "renote"
	MkNotificationTypeReply                 MkNotificationType = "reply"
	MkNotificationTypeMention               MkNotificationType = "mention"

	MkNotificationTypeUnknown MkNotificationType = "unknown"
)

type MkNotification struct {
	Id          string             `json:"id"`
	Type        MkNotificationType `json:"type"`
	UserId      *string            `json:"userId,omitempty"`
	CreatedAt   string             `json:"createdAt"`
	User        *MkUser            `json:"user,omitempty"`
	Note        *MkNote            `json:"note,omitempty"`
	Reaction    *string            `json:"reaction,omitempty"`
	Choice      *int               `json:"choice,omitempty"`
	Invitation  any                `json:"invitation"`
	Body        any                `json:"body"`
	Header      *string            `json:"header,omitempty"`
	Icon        *string            `json:"icon,omitempty"`
	Achievement string             `json:"achievement"`
}

func (n MkNotification) ToNotification(server string) (Notification, error) {
	r := Notification{
		Id:        n.Id,
		Type:      n.Type.ToNotificationType(),
		CreatedAt: n.CreatedAt,
	}
	var err error
	if n.User != nil {
		r.Account, err = n.User.ToAccount(server)
		if err != nil {
			return r, err
		}
	}
	if n.Note != nil {
		status := n.Note.ToStatus(server)
		r.Status = &status
	}
	return r, err
}

func (t MkNotificationType) ToNotificationType() NotificationType {
	switch t {
	case MkNotificationTypeNote:
		return NotificationTypeStatus
	case MkNotificationTypeFollow:
		return NotificationTypeFollow
	case MkNotificationTypeReceiveFollowRequest:
		return NotificationTypeFollowRequest
	case MkNotificationTypeReceiveReaction:
		return NotificationTypeFavourite
	case MkNotificationTypeReceiveRenote:
		return NotificationTypeReblog
	case MkNotificationTypeReply, MkNotificationTypeMention:
		return NotificationTypeMention
	default:
		return NotificationTypeUnknown
	}
}
