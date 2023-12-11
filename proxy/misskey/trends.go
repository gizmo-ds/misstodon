package misskey

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
)

func TrendsTags(ctx misstodon.Context, limit, offset int) ([]models.Tag, error) {
	var result []struct {
		Tag        string `json:"tag"`
		UsersCount int    `json:"usersCount"`
	}
	_, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(makeBody(ctx, nil)).
		SetResult(&result).
		Post("/api/hashtags/trend")
	if err != nil {
		return nil, err
	}
	var tags []models.Tag
	for _, r := range result {
		tag := models.Tag{
			Name: r.Tag,
			Url:  utils.JoinURL(*ctx.HOST(), "/tags/", r.Tag),
			History: []struct {
				Day      string `json:"day"`
				Uses     string `json:"uses"`
				Accounts string `json:"accounts"`
			}{
				{
					Day:      fmt.Sprint(time.Now().Unix()),
					Uses:     strconv.Itoa(r.UsersCount),
					Accounts: strconv.Itoa(r.UsersCount),
				},
			},
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func TrendsStatus(ctx misstodon.Context, limit, offset int) ([]models.Status, error) {
	var statuses []models.Status
	var result []models.MkNote
	_, err := client.R().
		SetBaseURL(ctx.ProxyServer()).
		SetBody(makeBody(ctx, utils.Map{"limit": limit})).
		SetResult(&result).
		Post("/api/notes/featured")
	if err != nil {
		return nil, err
	}
	for _, note := range result {
		statuses = append(statuses, note.ToStatus(ctx))
	}
	return statuses, nil
}
