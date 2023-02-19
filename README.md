# misstodon (WIP)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gizmo-ds/misstodon?style=flat-square)
[![Build images](https://img.shields.io/github/actions/workflow/status/gizmo-ds/misstodon/images.yaml?branch=main&label=docker%20image&style=flat-square)](https://github.com/gizmo-ds/misstodon/actions/workflows/images.yaml)
[![Release](https://img.shields.io/github/v/release/gizmo-ds/misstodon.svg?include_prereleases&style=flat-square)](https://github.com/gizmo-ds/misstodon/releases/latest)
[![License](https://img.shields.io/github/license/gizmo-ds/misstodon?style=flat-square)](./LICENSE)

## Progress

<details>

- [ ] .well-known
  - [x] /.well-known/webfinger
  - [x] /.well-known/nodeinfo
- [ ] Nodeinfo
  - [x] /nodeinfo/2.0
- [ ] Auth
  - [x] /oauth/authorize
  - [x] /oauth/token
  - [x] /api/v1/apps
  - [ ] /api/v1/apps/verify_credentials
- [ ] Instance
  - [x] /api/v1/instance
- [ ] Accounts
  - [x] /api/v1/accounts/lookup
  - [x] /api/v1/accounts/verify_credentials
  - [ ] /api/v1/accounts/update_credentials
  - [ ] /api/v1/accounts/relationships
  - [ ] /api/v1/accounts/:user_id/statuses
  - [ ] /api/v1/accounts/:user_id/following
  - [ ] /api/v1/accounts/:user_id/followers
- [ ] Statuses
  - [x] /api/v1/statuses/:status_id
  - [ ] /api/v1/statuses/:status_id/context
  - [ ] /api/v1/statuses/:status_id/favourite
  - [ ] /api/v1/statuses/:status_id/bookmark
- [ ] Timelines
  - [ ] /api/v1/timelines/home
  - [x] /api/v1/timelines/public
- [ ] Favourites
  - [ ] /api/v1/favourites
- [ ] Bookmarks
  - [ ] /api/v1/bookmarks
- [ ] Push
  - [ ] /api/v1/notifications
- [ ] Streaming
  - [ ] /api/v1/streaming
- [ ] Search
  - [ ] /api/v2/search
- [ ] Conversations
  - [ ] /api/v1/conversations
- [ ] Trends
  - [ ] /api/v1/trends/statuses

</details>
