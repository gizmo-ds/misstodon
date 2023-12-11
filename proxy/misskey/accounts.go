package misskey

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func AccountsLookup(ctx Context, acct string) (models.Account, error) {
	var host *string
	var info models.Account
	username, _host := utils.AcctInfo(acct)
	if username == "" {
		return info, ErrAcctIsInvalid
	}
	if _host == "" {
		_host = ctx.ProxyServer()
	}
	if _host != ctx.ProxyServer() {
		host = &_host
	}
	var result models.MkUser
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(map[string]any{
			"username": username,
			"host":     host,
		}).
		SetResult(&result).
		Post("/api/users/show")
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return info, ErrNotFound
	}
	return result.ToAccount(ctx.ProxyServer())
}

func AccountsStatuses(
	ctx Context, uid string, limit int,
	pinnedOnly, onlyMedia, onlyPublic, excludeReplies, excludeReblogs bool,
	maxID, minID string) ([]models.Status, error) {
	var notes []models.MkNote
	body := makeBody(ctx, utils.Map{
		"userId":         uid,
		"limit":          limit,
		"includeReplies": !excludeReplies,
	})
	if onlyMedia {
		body["fileType"] = SupportedMimeTypes
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&notes).
		Post("/api/users/notes")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("failed to get statuses")
	}
	statuses := slice.Map(notes, func(_i int, note models.MkNote) models.Status { return note.ToStatus(ctx.ProxyServer()) })
	return statuses, nil
}

func VerifyCredentials(ctx Context) (models.CredentialAccount, error) {
	var account models.Account
	var result models.MkUser
	var info models.CredentialAccount
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(utils.Map{"i": ctx.Token()}).
		SetResult(&result).
		Post("/api/i")
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return info, errors.New("failed to verify credentials")
	}
	account, err = result.ToAccount(ctx.ProxyServer())
	if err != nil {
		return info, err
	}
	info.Account = account
	if result.Description != nil {
		info.Source.Note = *result.Description
	}
	info.Source.Fields = info.Account.Fields
	return info, nil
}

// UpdateCredentials updates the credentials of the user.
func UpdateCredentials(ctx Context,
	displayName, note *string,
	locked, bot, discoverable *bool,
	sourcePrivacy *string, sourceSensitive *bool, sourceLanguage *string,
	fields []models.AccountField,
	avatar, header *multipart.FileHeader,
) (models.CredentialAccount, error) {
	var info models.CredentialAccount
	var body = utils.Map{
		"i": ctx.Token(),
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
		avatarFile, err := driveFileCreate(ctx, avatar.Filename, file)
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
		headerFile, err := driveFileCreate(ctx, header.Filename, file)
		if err != nil {
			return info, errors.WithStack(err)
		}
		body["bannerId"] = headerFile.ID
	}

	var result models.MkUser
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Patch("/api/i/update")
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return info, errors.New("failed to verify credentials")
	}
	account, err := result.ToAccount(ctx.ProxyServer())
	if err != nil {
		return info, err
	}
	info.Account = account
	if result.Description != nil {
		info.Source.Note = *result.Description
	}
	info.Source.Fields = account.Fields
	return info, err
}

func AccountFollowRequests(ctx Context,
	limit int, sinceID, maxID string) ([]models.Account, error) {
	var result []struct {
		ID       string        `json:"id"`
		Follower models.MkUser `json:"follower"`
		Followee models.MkUser `json:"followee"`
	}
	body := utils.Map{"i": ctx.Token(), "limit": limit}
	if sinceID != "" {
		body["sinceId"] = sinceID
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/following/requests/list")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var accounts []models.Account
	for _, r := range result {
		if a, err := r.Follower.ToAccount(ctx.ProxyServer()); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}

func AccountFollowRequestsCancel(ctx Context, userID string) error {
	data := utils.Map{"i": ctx.Token(), "userId": userID}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(data).
		Post("/api/following/requests/cancel")
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountFollowRequestsAccept(ctx Context, userID string) error {
	data := utils.Map{"i": ctx.Token(), "userId": userID}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(data).
		Post("/api/following/requests/accept")
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountFollowRequestsReject(ctx Context, userID string) error {
	data := utils.Map{"i": ctx.Token(), "userId": userID}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(data).
		Post("/api/following/requests/reject")
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountFollowers(ctx Context, userID string,
	limit int, sinceID, minID, maxID string) ([]models.Account, error) {
	var result []struct {
		ID         string        `json:"id"`
		CreatedAt  string        `json:"createdAt"`
		FolloweeId string        `json:"followeeId"`
		FollowerId string        `json:"followerId"`
		Follower   models.MkUser `json:"follower"`
	}
	body := makeBody(ctx, utils.Map{"limit": limit, "userId": userID})
	if v, ok := utils.StrEvaluation(sinceID, minID); ok {
		body["sinceId"] = v
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/users/followers")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}

	var accounts []models.Account
	for _, r := range result {
		if a, err := r.Follower.ToAccount(ctx.ProxyServer()); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}

func AccountFollowing(ctx Context,
	userID string,
	limit int, sinceID, minID, maxID string) ([]models.Account, error) {
	var result []struct {
		ID         string        `json:"id"`
		CreatedAt  string        `json:"createdAt"`
		FolloweeId string        `json:"followeeId"`
		FollowerId string        `json:"followerId"`
		Followee   models.MkUser `json:"followee"`
	}
	body := makeBody(ctx, utils.Map{"limit": limit, "userId": userID})
	if v, ok := utils.StrEvaluation(sinceID, minID); ok {
		body["sinceId"] = v
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/users/following")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}

	var accounts []models.Account
	for _, r := range result {
		if a, err := r.Followee.ToAccount(ctx.ProxyServer()); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}

func AccountRelationships(ctx Context,
	userIDs []string) ([]models.Relationship, error) {
	data := utils.Map{"i": ctx.Token(), "userId": userIDs}
	var result []models.MkRelation
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(data).
		SetResult(&result).
		Post("/api/users/relation")
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

func AccountFollow(ctx Context, userID string) error {
	data := utils.Map{"i": ctx.Token(), "userId": userID}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(data).
		Post("/api/following/create")
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountUnfollow(ctx Context, userID string) error {
	data := makeBody(ctx, utils.Map{"userId": userID})
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(data).
		Post("/api/following/delete")
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountMute(ctx Context, userID string, duration int64) error {
	body := makeBody(ctx, utils.Map{"userId": userID})
	if duration > 0 {
		body["expiresAt"] = time.Now().Add(time.Second * time.Duration(duration)).UnixMilli()
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		Post("/api/mute/create")
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountUnmute(ctx Context, userID string) error {
	data := makeBody(ctx, utils.Map{"userId": userID})
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(data).
		Post("/api/mute/delete")
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountsGet(ctx Context, userID string) (models.Account, error) {
	var info models.Account
	var result models.MkUser
	body := makeBody(ctx, utils.Map{"userId": userID})
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/users/show")
	if err != nil {
		return info, errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return info, ErrNotFound
	}
	return result.ToAccount(ctx.ProxyServer())
}

func AccountFavourites(ctx Context,
	limit int, sinceID, minID, maxID string,
) ([]models.Status, error) {
	type reactionsResult struct {
		ID        string        `json:"id"`
		User      models.MkUser `json:"user"`
		Note      models.MkNote `json:"note"`
		Type      string        `json:"type"`
		CreatedAt time.Time     `json:"createdAt"`
	}
	var result []reactionsResult
	body := makeBody(ctx, utils.Map{"limit": limit, "userId": ctx.UserID()})
	if v, ok := utils.StrEvaluation(sinceID, minID); ok {
		body["sinceId"] = v
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(body).
		SetResult(&result).
		Post("/api/users/reactions")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	return slice.Map(result, func(_ int, r reactionsResult) models.Status { return r.Note.ToStatus(ctx.ProxyServer()) }), nil
}
