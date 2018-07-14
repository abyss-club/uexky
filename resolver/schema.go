package resolver

// Schema for api
const Schema = `
schema {
    query: Query
    mutation: Mutation
}

type Query {
    profile(): User!
    threadSlice(tags: [String!], query: SliceQuery!): ThreadSlice!
    thread(id: String!): Thread!
    post(id: String!): Post!
    tags(query: String): TagTree!
}

type Mutation {
    auth(email: String!): Boolean!
    setName(name: String!): User!
    syncTags(tags: [String]!): User!
    pubThread(thread: ThreadInput!): Thread!
    pubPost(post: PostInput!): Post!
}

type SliceInfo {
    firstCursor: String!
    lastCursor: String!
}

input SliceQuery {
    before: String
    after: String
    limit: Int!
}

scalar Time

// Data Type Defines
type User {
    email: String!
    name: String
    tags: [String!]
}

input ThreadInput {
    anonymous: Boolean!
    content: String!
    mainTag: String!
    subTags: [String!]
    title: String
}

type Thread {
    id: String!
    anonymous: Boolean!
    content: String!
    createTime: Time!

    mainTag: String!
    subTags: [String!]
    title: String
    replies(query: SliceQuery!): PostSlice!
}

type ThreadSlice {
  threads: [Thread]!
  sliceInfo: SliceInfo!
}


input PostInput {
    threadID: String!
    anonymous: Boolean!
    content: String!
    refers: [String!]
}

type Post {
    id: String!
    anonymous: Boolean!
    content: String!
    createTime: Time!
    refers: [Post!]
}

type PostSlice {
  posts: [Post]!
  sliceInfo: SliceInfo!
}

type Tags {
    mainTags: [String!]!
    recommend: [String!]!
    tree: [TagTreeNode!]
}

type TagTreeNode {
    mainTag: String!
    subTags: [String!]
}
`
