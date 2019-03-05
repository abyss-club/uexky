export default `
  type Query {
    # Return a config object.
    config: Config!
  }

  type Mutation {
    editConfig(config: ConfigInput!): Config!
  }

  type Config {
    rateLimit: RateLimit!
    rateCost: RateCost!
  }

  type RateLimit {
    httpHeader: String!
    queryLimit: Int!
    queryResetTime: Int!
    mutLimit: Int!
    mutResetTime: Int!
  }

  type RateCost {
    createUser: Int!
    pubThread: Int!
    pubPost: Int!
  }

  input ConfigInput {
    rateLimit: RateLimitInput
    rateCost: RateCostInput
  }

  input RateLimitInput {
    httpHeader: String
    queryLimit: Int
    queryResetTime: Int
    mutLimit: Int
    mutResetTime: Int
  }

  input RateCostInput {
    createUser: Int
    pubThread: Int
    pubPost: Int
  }
`;
