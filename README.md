# misstodon

[![爱发电](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fafdian.net%2Fapi%2Fuser%2Fget-profile%3Fuser_id%3D75e549844b5111ed8df552540025c377&query=%24.data.user.name&label=%E7%88%B1%E5%8F%91%E7%94%B5&color=%23946ce6)](https://afdian.net/a/gizmo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gizmo-ds/misstodon?style=flat-square)
[![Build images](https://img.shields.io/github/actions/workflow/status/gizmo-ds/misstodon/images.yaml?branch=main&label=docker%20image&style=flat-square)](https://github.com/gizmo-ds/misstodon/actions/workflows/images.yaml)
[![License](https://img.shields.io/github/license/gizmo-ds/misstodon?style=flat-square)](./LICENSE)

Misskey Mastodon-compatible APIs, Getting my [Misskey](https://github.com/misskey-dev/misskey/tree/13.2.0) instance to work in [Elk](https://github.com/elk-zone/elk)

> **Warning**  
> This project is still in the early stage of development, and is not ready for production use.

## Demo

Elk: [https://elk.zone/misstodon.aika.dev/public](https://elk.zone/misstodon.aika.dev/public)  
Elk: [https://elk.zone/mt_misskey_moe.aika.dev/public](https://elk.zone/mt_misskey_moe.aika.dev/public)  
Phanpy: [https://phanpy.social/#/mt_misskey_io.aika.dev/p](https://phanpy.social/#/mt_misskey_io.aika.dev/p)

## How to Use

> **Warning**  
> `aika.dev` is a demonstration site and may not guarantee high availability. We recommend [self-hosting](#running-your-own-instance) for greater control.

### Domain Name Prefixing Scheme (Recommended)

The simplest usage method is to specify the instance address using a domain name prefix.

1. Replace underscores ("\_") in the domain name with double underscores ("\_\_").
2. Replace dots (".") in the domain name with underscores ("\_").
3. Prepend "mt\_" to the modified string.
4. Append ".aika.dev" to the modified string.

When processing `misskey.io` according to the described steps, it will be transformed into the result: `mt_misskey_io.aika.dev`.

```bash
curl --request GET --url 'https://mt_misskey_io.aika.dev/nodeinfo/2.0' | jq .
```

### Self-Hosting with Default Instance Configuration

Edit the 'config.toml' file, and within the "[proxy]" section, modify the "fallback_server" field. For example:

```toml
[proxy]
fallback_server = "misskey.io"
```

If you are [deploying using Docker Compose](#running-your-own-instance), you can specify the default instance by modifying the 'docker-compose.yml' file. Look for the 'MISSTODON_FALLBACK_SERVER' field within the Docker Compose configuration and set it to the desired default instance.

### Instance Specification via Query Parameter

```bash
curl --request GET --url 'https://misstodon.aika.dev/nodeinfo/2.0?server=misskey.io' | jq .
```

### Instance Specification via Header

```bash
curl --request GET --url https://misstodon.aika.dev/nodeinfo/2.0 --header 'x-proxy-server: misskey.io' | jq .
```

## Running your own instance

The simplest way is to use Docker Compose. Download the [docker-compose.yml](https://github.com/gizmo-ds/misstodon/raw/main/docker-compose.yml) file to your local machine. Customize it to your needs, particularly by changing the "MISSTODON_FALLBACK_SERVER" in the "environment" to your preferred Misskey instance domain. Afterward, run the following command:

```bash
docker-compose up -d
```

> **Important**  
> For security and privacy, we strongly discourage using HTTP directly. Instead, consider configuring a TLS certificate or utilizing Misstodon's AutoTLS feature for enhanced security.

## Roadmap

<details>

- [x] .well-known
  - [x] `GET` /.well-known/host-meta
  - [x] `GET` /.well-known/webfinger
  - [x] `GET` /.well-known/nodeinfo
  - [x] `GET` /.well-known/change-password
- [x] Nodeinfo
  - [x] `GET` /nodeinfo/2.0
- [ ] Auth
  - [x] `GET` /oauth/authorize
  - [x] `POST` /oauth/token
  - [x] `POST` /api/v1/apps
  - [ ] `GET` /api/v1/apps/verify_credentials
- [x] Instance
  - [x] `GET` /api/v1/instance
  - [x] `GET` /api/v1/custom_emojis
- [ ] Accounts
  - [x] `GET` /api/v1/accounts/lookup
  - [x] `GET` /api/v1/accounts/:user_id
  - [x] `GET` /api/v1/accounts/verify_credentials
  - [ ] `PATCH` /api/v1/accounts/update_credentials
  - [x] `GET` /api/v1/accounts/relationships
  - [ ] `GET` /api/v1/accounts/:user_id/statuses
  - [x] `GET` /api/v1/accounts/:user_id/following
  - [x] `GET` /api/v1/accounts/:user_id/followers
  - [x] `POST` /api/v1/accounts/:user_id/follow
  - [x] `POST` /api/v1/accounts/:user_id/unfollow
  - [x] `GET` /api/v1/follow_requests
  - [x] `POST` /api/v1/accounts/:user_id/mute
  - [x] `POST` /api/v1/accounts/:user_id/unmute
  - [x] `GET` /api/v1/bookmarks
  - [x] `GET` /api/v1/favourites
  - [ ] `GET` /api/v1/preferences
- [ ] Statuses
  - [x] `POST` /api/v1/statuses
  - [x] `GET` /api/v1/statuses/:status_id
  - [x] `DELETE` /api/v1/statuses/:status_id
  - [x] `GET` /api/v1/statuses/:status_id/context
  - [x] `POST` /api/v1/statuses/:status_id/reblog
  - [x] `POST` /api/v1/statuses/:status_id/favourite
  - [x] `POST` /api/v1/statuses/:status_id/unfavourite
  - [x] `POST` /api/v1/statuses/:status_id/bookmark
  - [x] `POST` /api/v1/statuses/:status_id/unbookmark
  - [ ] `GET` /api/v1/statuses/:status_id/favourited_by
  - [ ] `GET` /api/v1/statuses/:status_id/reblogged_by
- [x] Timelines
  - [x] `GET` /api/v1/timelines/home
  - [x] `GET` /api/v1/timelines/public
  - [x] `GET` /api/v1/timelines/tag/:hashtag
  - [ ] `WS` /api/v1/streaming
- [ ] Notifications
  - [x] `GET` /api/v1/notifications
  - [ ] `POST` /api/v1/push/subscription
  - [ ] `GET` /api/v1/push/subscription
  - [ ] `PUT` /api/v1/push/subscription
  - [ ] `DELETE` /api/v1/push/subscription
- [ ] Search
  - [ ] `GET` /api/v2/search
- [ ] Conversations
  - [ ] `GET` /api/v1/conversations
  - [ ] `DELETE` /api/v1/conversations/:id
  - [ ] `POST` /api/v1/conversations/:id/read
- [x] Trends
  - [x] `GET` /api/v1/trends/statuses
  - [x] `GET` /api/v1/trends/tags
- [x] Media
  - [x] `POST` /api/v1/media
  - [x] `POST` /api/v2/media

</details>

## Information for Developers

[Contributing](./CONTRIBUTING.md) Information about contributing to this project.

## Sponsors

[![Sponsors](https://afdian-connect.deno.dev/sponsor.svg)](https://afdian.net/a/gizmo)

## Contributors

![Contributors](https://contributors.aika.dev/gizmo-ds/misstodon/contributors.svg?align=left)
