package misskey

import (
	"mime/multipart"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAcctIsInvalid = errors.New("acct format is invalid")
)

func Lookup(server string, acct string) (models.Account, error) {
	var host *string
	var info models.Account
	username, _host := utils.AcctInfo(acct)
	if username == "" {
		return info, ErrAcctIsInvalid
	}
	if _host == "" {
		_host = server
	}
	if _host != server {
		host = &_host
	}
	var serverInfo models.MkUser
	resp, err := client.R().
		SetBody(map[string]any{
			"username": username,
			"host":     host,
		}).
		SetResult(&serverInfo).
		Post("https://" + server + "/api/users/show")
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return info, ErrNotFound
	}
	return serverInfo.ToAccount(server)
}

func AccountsStatuses(
	server, accountId string,
	limit int,
	pinnedOnly, onlyMedia, onlyPublic, excludeReplies, excludeReblogs bool,
	maxID, minID string) ([]models.Status, error) {
	var notes []models.MkNote
	r := map[string]any{
		"userId":         accountId,
		"limit":          limit,
		"includeReplies": !excludeReplies,
	}
	if onlyMedia {
		r["fileType"] = SupportedMimeTypes
	}
	resp, err := client.R().
		SetBody(r).
		SetResult(&notes).
		Post("https://" + server + "/api/users/notes")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("failed to get statuses")
	}
	var statuses []models.Status
	for _, note := range notes {
		statuses = append(statuses, note.ToStatus(server))
	}
	return statuses, nil
}

func VerifyCredentials(server, token string) (models.CredentialAccount, error) {
	var account models.Account
	var serverInfo models.MkUser
	var info models.CredentialAccount
	resp, err := client.R().
		SetBody(utils.Map{"i": token}).
		SetResult(&serverInfo).
		Post("https://" + server + "/api/i")
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return info, errors.New("failed to verify credentials")
	}
	account, err = serverInfo.ToAccount(server)
	if err != nil {
		return info, err
	}
	info.Account = account
	if serverInfo.Description != nil {
		info.Source.Note = *serverInfo.Description
	}
	info.Source.Fields = info.Account.Fields
	return info, nil
}

// UpdateCredentials updates the credentials of the user.
func UpdateCredentials(server, token string,
	displayName, note *string,
	locked, bot, discoverable *bool,
	sourcePrivacy *string, sourceSensitive *bool, sourceLanguage *string,
	fields []models.AccountField,
// FIXME: ????????????????????????????????????
	avatar, header *multipart.FileHeader,
) (models.CredentialAccount, error) {
	var info models.CredentialAccount
	var body = utils.Map{
		"i": token,
	}
	if displayName != nil {
		body["name"] = displayName
	}
	if note != nil {
		body["description"] = note
	}
	if locked != nil {
		body["isLocked"] = locked
	}
	if bot != nil {
		body["isBot"] = bot
	}
	if sourceLanguage != nil {
		body["lang"] = sourceLanguage
	}
	if fields != nil {
		body["fields"] = fields
	}
	var serverInfo models.MkUser
	resp, err := client.R().
		SetBody(body).
		SetResult(&serverInfo).
		Patch("https://" + server + "/api/i/update")
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return info, errors.New("failed to verify credentials")
	}
	account, err := serverInfo.ToAccount(server)
	if err != nil {
		return info, err
	}
	info.Account = account
	if serverInfo.Description != nil {
		info.Source.Note = *serverInfo.Description
	}
	info.Source.Fields = account.Fields
	return info, err
}
