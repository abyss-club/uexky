export default `
  extend type Query {
    # A slice of thread.
    threadSlice(tags: [String!], query: SliceQuery!): ThreadSlice!
    # A thread object.
    thread(id: String!): Thread!
  }

  extend type Mutation {
    # Publish a new thread.
    pubThread(thread: ThreadInput!): Thread!
  }

  # Construct a new thread.
  input ThreadInput {
    # Toggle anonymousness. If true, a new ID will be generated in each thread.
    anonymous: Boolean!
    content: String!
    # Required. Only one mainTag is allowed.
    mainTag: String!
    # Optional, maximum of 4.
    subTags: [String!]
    # Optional. If not set, the title will be '无题'.
    title: String
  }

  type Thread {
    # UUID with 8 chars in length, and will increase to 9 after 30 years.
    id: String!
    createdAt: Time!
    # Thread was published anonymously or not.
    anonymous: Boolean!
    # Same format as id if anonymous, name of User otherwise.
    author: String!
    # Default to '无题'.
    title: String
    content: String!
    # Only one mainTag is allowed.
    mainTag: String!
    # Optional, maximum of 4.
    subTags: [String!]
    replies(query: SliceQuery!): PostSlice!
    replyCount: Int!
    catalog: [ThreadCatalogItem!]
    blocked: Boolean!
    locked: Boolean!
  }

  type ThreadCatalogItem {
    postId: [String!]
    createdAt: Time!
  }

  type ThreadSlice {
    threads: [Thread]!
    sliceInfo: SliceInfo!
  }
`;
