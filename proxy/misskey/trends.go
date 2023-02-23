package misskey

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
)

func TrendsTags(server, token string, limit, offset int) ([]models.Tag, error) {
	var result []struct {
		Tag        string `json:"tag"`
		UsersCount int    `json:"usersCount"`
	}
	_, err := client.R().
		SetBody(utils.Map{
			"i": token,
		}).
		SetResult(&result).
		Post("https://" + server + "/api/hashtags/trend")
	if err != nil {
		return nil, err
	}
	var tags []models.Tag
	for _, r := range result {
		tag := models.Tag{
			Name: r.Tag,
			Url:  "https://" + server + "/tags/" + r.Tag,
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
