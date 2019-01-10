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
    addSubbedTags(tags: [String!]!): User!
    # Delete tags subscribed by user.
    delSubbedTags(tags: [String!]!): User!

    # admin's apis:
    banUser(postId: String!)
    blockPost(postId: String!)
    lockThread(threadId: String!)
    blockThread(threadId: String!)
    editTags(threadId: String!, mainTag: String!, subTags: [String!]!)
  }

  type User {
    email: String!
    # The Name of user. Required when not posting anonymously.
    name: String
    # Tags saved by user.
    tags: [String!]
    role: UserRole
  }

  # user roles:
  # role: 'SuperAdmin': modify config, manage all tags.
  #       'TagAdmin':   manage several mainTags.
  #       None:         normal user. (won't have "role" field)
  type UserRole {
    role: String!
    params: [String!]
  }
`;
