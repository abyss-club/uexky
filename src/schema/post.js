export default `
  extend type Query {
    # A post object.
    post(id: String!): Post!
  }

  extend type Mutation {
    # Publish a new post.
    pubPost(post: PostInput!): Post!
  }

  # Input object describing a Post to be published.
  input PostInput {
    # ID of the replying thread's.
    threadId: String!
    # Should be sent as anonymous or not.
    anonymous: Boolean!
    # Markdown formatted content.
    content: String!
    # Set quoting PostIDs.
    quoteIds: [String!]
  }

  # Object describing a Post.
  type Post {
    id: String!
    createdAt: Time!
    anonymous: Boolean!
    # Name if not anonymous, anonymous ID otherwise.
    author: String!
    # Markdown formatted content.
    content: String!
    # Other posts that the post has quoted.
    quotes: [Post!]
    # Amount of times that the post is quoted.
    quotedCount: Int!
    # The post is blocked or not.
    blocked: Boolean!
  }

  # PostSlice object is for selecting specific 'slice' of Post objects to
  # return. Affects the returning SliceInfo.
  type PostSlice {
    posts: [Post]!
    sliceInfo: SliceInfo!
  }
`;
