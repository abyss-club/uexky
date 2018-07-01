# Uexky

一个新的讨论版程序。此为后端项目。

# 功能设计

## 总体

* 响应式设计
* 减少页面刷新和载入
* (TODO) 兼容PWA

## 用户

### 注册

* 发言的用户必须是注册用户，但是无需使用邮箱密码等方式，而是自动生成一个用户 token
* 用户 token 可以用于在其他设备登录
* (TODO) 如有需要，则可以添加用户名和密码登录。token 仍然保留，用以忘记密码时重置密码

### 发帖身份

* 用户发帖时可以选择匿名还是具名发帖
* 具名发帖需要设置发言ID，发言ID唯一且不可变
* (TODO) 一个用户可以拥有多个发言ID

### 匿名发帖

* 用户选择匿名发帖后，会自动生成一个匿名ID代表用户
* 同一个用户在同一个贴内的匿名ID相同，不同的帖子里面不同
* (TODO) 匿名ID可以重置，重置之后，在同一贴发帖，匿名ID也不会相同

## 内容

### 主帖(Thread) 和标签 (Tag)

* 主贴分类使用 标签/Tag 的方式，必须有且仅有一个主 Tag，可以有 0-4 个子 Tag
* 主 Tag 只能从现有主 Tag 中选择，子 Tag 可以自主添加
* 用户可以订阅 Tag，可以在按已订阅的多个 Tag 合并在一个时间线中查看
* (TODO) 每个 Tag 有自己的专页保存内容，类似 Wiki 的形式，因此无需置顶功能
* (TODO) 用户可以提议合并 Tag，提升子 Tag 为主 Tag
* (TODO) 支持多个订阅组
* 主贴无需含有标题
* 其余与普通帖子相同

### 帖子（Post)

* 使用 Markdown 的格式发帖，支持 Markdown 内建的链接、图片等内容。不允许使用 html 代码添加自定义样式
* (TODO) 支持主流视频网站的内嵌播放器
* (TODO) 帖子内容可以编辑并可查看编辑历史

### 引用回复（Refer)

* 可以对主贴外的帖子进行回复（普通回帖即为针对主贴的回复）, 发表后会在回复内容上附上被回复内容的部分引用
* 回复可以针对多人
* 查看回复交互参考 Telegram，点击被引用的回复可以滚动到被回复的帖子，并提供一个按钮返回。在同一页面内完成，网页不应该加载。
* (TODO) 以会话的形式查看多人多次相互回复
* 在页面中提醒用户被回复，并可以在用户中心查看
* (TODO) 使用浏览器通知回复消息

## 管理员

* 屏蔽主贴，回帖。沉贴
* 封禁用户或者 IP
* (TODO) 修改帖子 Tag
* (TODO) 删除，合并，提升 Tag
* (TODO) 在帖子中附上管理员批注
* 解除上述的操作

# API

原计划的 RESTful 风格的 API 设计：

[API Reference](https://github.com/CrowsT/uexky/wiki/API-Reference)

现在打算使用更好的 GraphQL API 来设计，功能上和之前的 RESTful API 保持一致，Schema 请见：

[GraphQL Schema](https://github.com/CrowsT/uexky/blob/master/api/schema.gql)

# Build Instructions

Dependencies needed:
* `go`
* `dep`
* `redis`
* `mongodb`

### Redis & Mongo
Check your db names in `mongo` to ensure no database named `develop`.

Run `redis` and `mongod` before running Go server.

### Go
```bash
# Clone repo into GOPATH
$ go get gitlab.com/abyss.club/uexky

$ cd GOPATH/src/gitlab.com/abyss.club/uexky
$ cp config.sample.json config.json
# Install Go dependencies
$ dep ensure
# Compile and serve API server at port 5000
$ go run main.go -c config.json
# ... or build executable then run
$ go build && ./uexky
```

