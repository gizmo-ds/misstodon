package models

type MkStats struct {
	NotesCount         int `json:"notesCount"`
	UsersCount         int `json:"usersCount"`
	OriginalUsersCount int `json:"originalUsersCount"`
	OriginalNotesCount int `json:"originalNotesCount"`
	ReactionsCount     int `json:"reactionsCount"`
	Instances          int `json:"instances"`
	DriveUsageLocal    int `json:"driveUsageLocal"`
	DriveUsageRemote   int `json:"driveUsageRemote"`
}
