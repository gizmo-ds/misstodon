package models

type MkRelation struct {
	ID                             string `json:"id"`
	IsFollowing                    bool   `json:"isFollowing"`
	IsFollowed                     bool   `json:"isFollowed"`
	HasPendingFollowRequestFromYou bool   `json:"hasPendingFollowRequestFromYou"`
	HasPendingFollowRequestToYou   bool   `json:"hasPendingFollowRequestToYou"`
	IsBlocking                     bool   `json:"isBlocking"`
	IsBlocked                      bool   `json:"isBlocked"`
	IsMuted                        bool   `json:"isMuted"`
}

func (r MkRelation) ToRelationship() Relationship {
	return Relationship{
		ID:         r.ID,
		Following:  r.IsFollowing,
		FollowedBy: r.IsFollowed,
		Requested:  r.HasPendingFollowRequestFromYou,
		Languages:  []string{},
		Blocking:   r.IsBlocking,
		BlockedBy:  r.IsBlocked,
		Muting:     r.IsMuted,
	}
}
