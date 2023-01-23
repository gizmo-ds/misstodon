# misstodon (WIP)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gizmo-ds/misstodon?style=flat-square)
[![Release](https://img.shields.io/github/v/release/gizmo-ds/misstodon.svg?include_prereleases&style=flat-square)](https://github.com/gizmo-ds/misstodon/releases/latest)
[![License](https://img.shields.io/github/license/gizmo-ds/misstodon?style=flat-square)](./LICENSE)

## Progress

| Status             | API                                |
| ------------------ | ---------------------------------- |
| :white_check_mark: | /.well-known/webfinger             |
| :white_check_mark: | /.well-known/nodeinfo              |
| :white_check_mark: | /nodeinfo/2.0                      |
| :white_check_mark: | /oauth/authorize                   |
| :white_check_mark: | /oauth/token                       |
| :white_check_mark: | `v1` /instance                     |
| :white_check_mark: | `v1` /accounts/lookup              |
| :white_check_mark: | `v1` /accounts/verify_credentials  |
| :construction:     | `v1` /accounts/<user_id>/statuses  |
| :x:                | `v1` /accounts/<user_id>/following |
| :x:                | `v1` /accounts/<user_id>/followers |
| :construction:     | `v1` /statuses/<status_id>         |
| :construction:     | `v1` /statuses/<status_id>/context |
| :x:                | `v1` /notifications                |
| :x:                | `v1` /streaming                    |
| :x:                | `v2` /search                       |
| :white_check_mark: | `v1` /apps                         |
| :x:                | `v1` /accounts/relationships       |
| :construction:     | `v1` /timelines/home               |
| :construction:     | `v1` /timelines/public             |
| :question:         | `v1` /conversations                |
| :construction:     | `v1` /favourites                   |
| :x:                | `v1` /trends/statuses              |
