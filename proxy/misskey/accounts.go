package misskey

import (
	"mime/multipart"
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func AccountsLookup(server string, acct string) (models.Account, error) {
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
		Post(utils.JoinURL(server, "/api/users/show"))
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return info, ErrNotFound
	}
	return serverInfo.ToAccount(server)
}

func AccountsStatuses(
	server, uid, token string,
	limit int,
	pinnedOnly, onlyMedia, onlyPublic, excludeReplies, excludeReblogs bool,
	maxID, minID string) ([]models.Status, error) {
	var notes []models.MkNote
	r := map[string]any{
		"userId":         uid,
		"limit":          limit,
		"includeReplies": !excludeReplies,
	}
	if token != "" {
		r["i"] = token
	}
	if onlyMedia {
		r["fileType"] = SupportedMimeTypes
	}
	resp, err := client.R().
		SetBody(r).
		SetResult(&notes).
		Post(utils.JoinURL(server, "/api/users/notes"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("failed to get statuses")
	}
	statuses := lo.Map(notes, func(note models.MkNote, _i int) models.Status { return note.ToStatus(server) })
	return statuses, nil
}

func VerifyCredentials(server, token string) (models.CredentialAccount, error) {
	var account models.Account
	var serverInfo models.MkUser
	var info models.CredentialAccount
	resp, err := client.R().
		SetBody(utils.Map{"i": token}).
		SetResult(&serverInfo).
		Post(utils.JoinURL(server, "/api/i"))
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
	if avatar != nil {
		file, err := avatar.Open()
		if err != nil {
			return info, errors.WithStack(err)
		}
		defer file.Close()
		avatarFile, err := driveFileCreate(server, token, avatar.Filename, file)
		if err != nil {
			return info, errors.WithStack(err)
		}
		body["avatarId"] = avatarFile.ID
	}
	if header != nil {
		file, err := header.Open()
		if err != nil {
			return info, errors.WithStack(err)
		}
		defer file.Close()
		headerFile, err := driveFileCreate(server, token, header.Filename, file)
		if err != nil {
			return info, errors.WithStack(err)
		}
		body["bannerId"] = headerFile.ID
	}

	var serverInfo models.MkUser
	resp, err := client.R().
		SetBody(body).
		SetResult(&serverInfo).
		Patch(utils.JoinURL(server, "/api/i/update"))
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

func AccountFollowRequests(server, token string,
	limit int, sinceID, maxID string) ([]models.Account, error) {
	var result []struct {
		ID       string        `json:"id"`
		Follower models.MkUser `json:"follower"`
		Followee models.MkUser `json:"followee"`
	}
	body := utils.Map{"i": token, "limit": limit}
	if sinceID != "" {
		body["sinceId"] = sinceID
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(server, "/api/following/requests/list"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var accounts []models.Account
	for _, r := range result {
		if a, err := r.Follower.ToAccount(server); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}

func AccountFollowRequestsCancel(server, token string, accountID string) error {
	data := utils.Map{"i": token, "userId": accountID}
	resp, err := client.R().
		SetBody(data).
		Post(utils.JoinURL(server, "/api/following/requests/cancel"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountFollowRequestsAccept(server, token string, accountID string) error {
	data := utils.Map{"i": token, "userId": accountID}
	resp, err := client.R().
		SetBody(data).
		Post(utils.JoinURL(server, "/api/following/requests/accept"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountFollowRequestsReject(server, token string, accountID string) error {
	data := utils.Map{"i": token, "userId": accountID}
	resp, err := client.R().
		SetBody(data).
		Post(utils.JoinURL(server, "/api/following/requests/reject"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountFollowers(server, token string,
	accountID string,
	limit int, sinceID, minID, maxID string) ([]models.Account, error) {
	var result []struct {
		ID         string        `json:"id"`
		CreatedAt  string        `json:"createdAt"`
		FolloweeId string        `json:"followeeId"`
		FollowerId string        `json:"followerId"`
		Follower   models.MkUser `json:"follower"`
	}
	body := utils.Map{"limit": limit, "userId": accountID}
	if token != "" {
		body["i"] = token
	}
	if v, ok := utils.StrEvaluation(sinceID, minID); ok {
		body["sinceId"] = v
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(server, "/api/users/followers"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}

	var accounts []models.Account
	for _, r := range result {
		if a, err := r.Follower.ToAccount(server); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}

func AccountFollowing(server, token string,
	accountID string,
	limit int, sinceID, minID, maxID string) ([]models.Account, error) {
	var result []struct {
		ID         string        `json:"id"`
		CreatedAt  string        `json:"createdAt"`
		FolloweeId string        `json:"followeeId"`
		FollowerId string        `json:"followerId"`
		Followee   models.MkUser `json:"followee"`
	}
	body := utils.Map{"limit": limit, "userId": accountID}
	if token != "" {
		body["i"] = token
	}
	if v, ok := utils.StrEvaluation(sinceID, minID); ok {
		body["sinceId"] = v
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(server, "/api/users/following"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}

	var accounts []models.Account
	for _, r := range result {
		if a, err := r.Followee.ToAccount(server); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}

func AccountRelationships(server, token string,
	accountIDs []string) ([]models.Relationship, error) {
	data := utils.Map{"i": token, "userId": accountIDs}
	var result []models.MkRelation
	resp, err := client.R().
		SetBody(data).
		SetResult(&result).
		Post(utils.JoinURL(server, "/api/users/relation"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var relationships []models.Relationship
	for _, r := range result {
		relationships = append(relationships, r.ToRelationship())
	}
	return relationships, nil
}

func AccountFollow(server, token string, accountID string) error {
	data := utils.Map{"i": token, "userId": accountID}
	resp, err := client.R().
		SetBody(data).
		Post(utils.JoinURL(server, "/api/following/create"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountUnfollow(server, token string, accountID string) error {
	data := utils.Map{"i": token, "userId": accountID}
	resp, err := client.R().
		SetBody(data).
		Post(utils.JoinURL(server, "/api/following/delete"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountMute(server, token string, accountID string, expiresAt int64) error {
	data := utils.Map{"i": token, "userId": accountID}
	if expiresAt > 0 {
		data["expiresAt"] = expiresAt
	}
	resp, err := client.R().
		SetBody(data).
		Post(utils.JoinURL(server, "/api/mute/create"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountUnmute(server, token string, accountID string) error {
	data := utils.Map{"i": token, "userId": accountID}
	resp, err := client.R().
		SetBody(data).
		Post(utils.JoinURL(server, "/api/mute/delete"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountsGet(server, token, accountID string) (models.Account, error) {
	var info models.Account
	var serverInfo models.MkUser
	body := utils.Map{"userId": accountID}
	if token != "" {
		body["i"] = token
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&serverInfo).
		Post(utils.JoinURL(server, "/api/users/show"))
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return info, ErrNotFound
	}
	return serverInfo.ToAccount(server)
}
