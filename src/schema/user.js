/* tslint:disable:max-line-length */
const userQueries = `
  # A user profile object.
  profile(): User!
`;

const userMutations = `
  # Register/Login via email address. An email containing login info will be sent to the provided email address.
  auth(email: String!): Boolean!
  # Set the Name of user.
  setName(name: String!): User!
  # Directly edit tags subscribed by user.
  syncTags(tags: [String]!): User!
  # Add tags subscribed by user.
  addSubbedTags(tags: [String!]!): User!
  # Delete tags subscribed by user.
  delSubbedTags(tags: [String!]!): User!
`;

const typeDef = `
  extend type Query {
    # A user profile object.
    profile: User!
  }
  type User {
    email: String!
    # The Name of user. Required when not posting anonymously.
    name: String
    # Tags saved by user.
    tags: [String!]
  }
`;

export { userQueries, userMutations, typeDef };
