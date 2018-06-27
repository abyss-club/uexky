package api

const schema = `
schema {
    query: Query
    mutation: Mutation
}

type Query {
    account(): Account!
    threadSlice(limit: Int! tags: [String!] after: String!): ThreadSlice!
    thread(id: String!): Thread!
    post(id: String!): Post!
	uexky(): Uexky!
}

type Mutation {
	auth(email: String!): Boolean!
    addName(name: String!): Account!
    syncTags(tags: [String]!): Account!
    pubThread(thread: ThreadInput!): Thread!
    pubPost(post: PostInput!): Post!
}

type SliceInfo {
    firstCursor: String!
    lastCursor: String!
}

scalar Time

// Data Type Defines
type Account {
    token: String!
    names: [String!]
    tags: [String!]
}

input ThreadInput {
    author: String
    content: String!
    mainTag: String!
    subTags: [String!]
    title: String
}

type Thread {
    id: String!
    anonymous: Boolean!
    author: String!
    content: String!
    createTime: Time!

    mainTag: String!
    subTags: [String!]
    title: String
    replies(limit: Int! after: String before: String): PostSlice!
}

type ThreadSlice {
  threads: [Thread]!
  sliceInfo: SliceInfo!
}


input PostInput {
    threadID: String!
    author: String
    content: String!
    refers: [String!]
}

type Post {
    id: String!
    anonymous: Boolean!
    author: String!
    content: String!
    createTime: Time!
    refers: [Post!]
}

type PostSlice {
  posts: [Post]!
  sliceInfo: SliceInfo!
}

type Uexky {
	mainTags: [String!]!
}
`
