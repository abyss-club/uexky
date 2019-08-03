export default `
  extend type Query {
    # A user profile object.
    profile: User!
  }

  extend type Mutation {
    # Register/Login via email address.
    # An email containing login info will be sent to the provided email address.
    auth(email: String!): Boolean!
    # Set the Name of user.
    setName(name: String!): User!
    # Directly edit tags subscribed by user.
    syncTags(tags: [String]!): User!
    # Add tags subscribed by user.
    addSubbedTag(tag: String!): User!
    # Delete tags subscribed by user.
    delSubbedTag(tag: String!): User!

    # mod's apis:
    banUser(postId: String, threadId: String): Boolean!
    blockPost(postId: String!): Post!
    lockThread(threadId: String!): Thread!
    blockThread(threadId: String!): Thread!
    editTags(threadId: String!, mainTag: String!, subTags: [String!]!): Thread!
  }

  type User {
    email: String!
    # The Name of user. Required when not posting anonymously.
    name: String
    # Tags saved by user.
    tags: [String!]
    # user roles:
    # role: admin: modify config, manage mods.
    #       mod:   lock/block thread, lock/block post, ban user.
    #       null:  normal user.
    role: String

    # Threads published by the user.
    threads(query: SliceQuery!): ThreadSlice!
    # Threads replied by the user.
    posts(query: SliceQuery!): PostSlice!
  }
`;
