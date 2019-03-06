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
      threadId: String!
      anonymous: Boolean!
      content: String!
      # Set quoting PostIDs.
      quoteIds: [String!]
  }

  # Object describing a Post.
  type Post {
      id: String!
      anonymous: Boolean!
      author: String!
      content: String!
      createdAt: Time!
      quotes: [Post!]
      quoteCount: Int!
      blocked: Boolean!
  }

  # PostSlice object is for selecting specific 'slice' of Post objects to
  # return. Affects the returning SliceInfo.
  type PostSlice {
      posts: [Post]!
      sliceInfo: SliceInfo!
  }
`;
