# Uexky

讨论板 [abyss](https://gitlab.com/abyss.club/abyss) 的后端项目。

## 功能设计

参见 Abyss 的 [开发需求](https://gitlab.com/abyss.club/abyss#%E5%BC%80%E5%8F%91%E9%9C%80%E6%B1%82%E5%8F%8A%E7%9B%B8%E5%85%B3%E5%90%8D%E8%AF%8D)

## Build Instructions

Dependencies needed:
* `go`
* `dep`
* `redis`
* `mongodb`

### Redis & Mongo

Check your db names in `mongo` to ensure no database named `develop`.
Run `redis` and `mongod` before running Go server.

### Config File

Config file use json, structure is:

```go
var Config struct {
        Mongo struct {
                URI string `json:"mongo_uri"`
                DB  string `json:"db"`
        } `json:"mongo"`
        RedisURI string   `json:"redis_url"`
        MainTags []string `json:"main_tags"`
        Proto    string   `json:"proto"`
        Domain   struct {
                WEB string `json:"web"`
                API string `json:"api"`
        } `json:"domain"`
        Mail struct {
                Domain     string `json:"domain"`
                PrivateKey string `json:"private_key"`
                PublicKey  string `json:"public_key"`
        } `json:"mail"`
}
```

Default Value is:

```json
{
    "mongo": {
        "mongo_url": "localhost:27017",
        "db": "develop"
    },
    "redis_url": "redis://localhost:6379/0",
    "proto": "https",
    "domain": {
        "web": "abyss.club",
        "api": "api.abyss.club"
    },
    "mail": {
        "domain": "mail.abyss.club"
    }
}
```

### Go

```bash
# Clone repo into GOPATH
$ go get gitlab.com/abyss.club/uexky

$ cd $GOPATH/src/gitlab.com/abyss.club/uexky
$ cp config.sample.json config.json
# Install Go dependencies
$ dep ensure
# Compile and serve API server at port 5000
$ go run main.go -c config.json
# ... or build executable then run
$ go build && ./uexky -c config.json
```

## API 说明

API 使用 [graphql](https://graphql.org/)，API schema [见此](https://gitlab.com/abyss.club/abyss/blob/master/api.gql)

### 登录/注册流程

登录注册使用同一个流程，没有区分，以下简称登录。

1. 登录需要提供一个 email 地址：

```
type Mutation {
    auth(email: String!): Boolean!
}
```

2. 如果返回 `true`，则会往提供的地址发送邮件，邮件中包含一个登录链接。点击登录
链接将会被重定向至网站首页并设置好 cookie，此时登录完成。
