# Uexky [![pipeline status](https://gitlab.com/abyss.club/uexky/badges/master/pipeline.svg)](https://gitlab.com/abyss.club/uexky/commits/master) [![coverage report](https://gitlab.com/abyss.club/uexky/badges/master/coverage.svg)](https://gitlab.com/abyss.club/uexky/commits/master)

讨论板 [abyss](https://gitlab.com/abyss.club/abyss) 的后端项目。

## 功能设计

参见 Abyss 的 [开发需求](https://gitlab.com/abyss.club/abyss#%E5%BC%80%E5%8F%91%E9%9C%80%E6%B1%82%E5%8F%8A%E7%9B%B8%E5%85%B3%E5%90%8D%E8%AF%8D)

## Build Instructions

Dependencies needed:

* `go`
* `dep`
* `redis`
* `mongodb`

### Config File

Config file use json, structure is:

```go
var Config struct {
    Mongo struct {
        URL string `json:"url"`
        DB  string `json:"db"`
    } `json:"mongo"`
    RedisURL string   `json:"redis_url"`
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
    RateLimit struct {
        QueryLimit        int `json:"query_limit"`
        QueryResetTime    int `json:"query_reset_time"`
        MutLimit     int `json:"mut_limit"`
        MutResetTime int `json:"mut_reset_time"`
        Cost              struct {
            CreateUser int `json:"create_user"`
            PubThread  int `json:"pub_thread"`
            PubPost    int `json:"pub_post"`
        } `json:"cost"`
    } `json:"rate_limit"`
}
```

Default Value is:

```json
{
    "mongo": {
        "url": "localhost:27017",
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
    "rate_limit": {
        "query_limit": 300,
        "query_reset_time": 3600,
        "mut_limit": 30,
        "mut_reset_time": 3600,
        "pub_thread_count": 10,
        "pub_post_count": 1,
        "create_count": 30,
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

### 分页

基于 `游标/Cursor` 分页，游标为一个字符串。客户端无需理解游标的含义，
也不应该利用游标的含义编程。使用分页查询的接口，需提供类型 `SliceQuery` 的参数：

```
input SliceQuery {
    before: String
    after: String
    limit: Int!
}
```

其中，`before` 和 `after` 中必须指定一个，如果指定空字符串表示 最后一个之前/第一个之后。
返回值将会以如下方式展示（以 Thread 为例）：

```
type ThreadSlice {
    threads: [Thread]!
    sliceInfo: SliceInfo!
}
```

`threads` 即为查询所得的 thread 列表。注意，分页的查询将会附加返回类型为 `SliceInfo` 的参数：

```
type SliceInfo {
    firstCursor: String!
    lastCursor: String!
}
```

包含了代表返回值 threads 中第一个和最后一个对象的游标。可以基于此游标开始下一次查询。
