# qiita

Go client library for [qiita API v2](https://qiita.com/api/v2/docs).

[![CircleCI](https://circleci.com/gh/muiscript/qiita/tree/master.svg?style=svg)](https://circleci.com/gh/muiscript/qiita/tree/master)
[![codecov](https://codecov.io/gh/muiscript/qiita/branch/master/graph/badge.svg)](https://codecov.io/gh/muiscript/qiita)

## usage

```go
logger := log.New(os.Stdout, "[LOG]", log.LstdFlags)
qiita := qiita.New("<YOUR_ACCESS_TOKEN>", logger)

ctx := context.Background()

// get user
user, err := qiita.GetUser(ctx, "muiscript")

// get item
item, err := qiita.GetItem(ctx, "b4ca1773580317e7112e")
```

## API list

#### apis available for unauthorized/authorized users

|  | Endpoint | Method Signature |
| --- | --- | --- |
| :heavy_check_mark: | `GET` - `/users` | `GetUsers(ctx context.Context, page int, perPage int)` |
| :heavy_check_mark: | `GET` - `/users/:user_id` | `GetUser(ctx context.Context, userID string)` |
| :heavy_check_mark: | `GET` - `/users/:user_id/followees` | `GetFollowees(ctx context.Context, userID string, page int, perPage int)` |
|  | `GET` - `/users/:user_id/followers` | |
|  | `GET` - `/users/:user_id/items` | |
|  | `GET` - `/users/:user_id/stocks` | |
|  | `GET` - `/users/:user_id/following_tags` | |
|  | `GET` - `/items` | |
| :heavy_check_mark: | `GET` - `/items/:item_id` | `GetItem(ctx context.Context, itemID string)` |
|  | `GET` - `/items/:item_id/stockers` | |
|  | `GET` - `/items/:item_id/comments` | |
|  | `GET` - `/tags` | |
|  | `GET` - `/tags/:tag_id` | |
|  | `GET` - `/tags/:tag_id/items` | |
|  | `GET` - `/comments/:comment_id` | |

#### apis available for unauthorized/authorized users

|  | Endpoint | Method Signature |
| --- | --- | --- |
| :heavy_check_mark: | `GET` - `/users/:user_id/following` | `IsFollowingUser(ctx context.Context, userID string)` |
|  | `DELETE` - `/users/:user_id/following` | |
|  | `PUT` - `/users/:user_id/following` | |
|  | `GET` - `/authenticated_user` | |
|  | `GET` - `/authenticated_user/items` | |
|  | `POST` - `/items` | |
|  | `DELETE` - `/items/:item_id` | |
|  | `PATCH` - `/items/:item_id` | |
|  | `DELETE` - `/items/:item_id` | |
|  | `GET` - `/items/:item_id/stock` | |
|  | `PUT` - `/items/:item_id/stock` | |
|  | `DELETE` - `/items/:item_id/stock` | |
|  | `GET` - `/tags/:tag_id/following` | |
|  | `PUT` - `/tags/:tag_id/following` | |
|  | `DELETE` - `/tags/:tag_id/following` | |
|  | `POST` - `/items/:item_id/comments` | |
|  | `PATCH` - `/comments/:comment_id` | |
|  | `DELETE` - `/comments/:comment_id` | |
|  | `PUT` - `/comments/:comment_id/thank` | |
|  | `DELETE` - `/comments/:comment_id/thank` | |
