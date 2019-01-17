export default `
  type Query {
    # Return a config object.
    config: Config!
  }

  type Mutation {
    editConfig(config: ConfigInput!): Config!
  }

  type Config {
    # Main tags are mandatory and predefined manually.
    mainTags: [String!]!
    # RateLimit is a JSON string.
    rateLimit: String!
  }

  input ConfigInput {
    # Update Main tags.
    mainTags: [String!]
    # Provide a JSON string to modify rateLimit.
    rateLimit: String
  }
`;
