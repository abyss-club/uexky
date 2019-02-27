export default `
  type Query {
    # Return a config object.
    config: Config!
  }

  type Mutation {
    editConfig(config: ConfigInput!): Config!
  }

  type Config {
    rateLimit: ReteLimit!
    cost: RateCost!
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
    rateLimit: ReteLimitInput
    cost: RateCostLimitInput
  }

  type RateLimitInput {
    httpHeader: String
    queryLimit: Int
    queryResetTime: Int
    mutLimit: Int
    mutResetTime: Int
  }

  type RateCostInput {
    createUser: Int
    pubThread: Int
    pubPost: Int
  }
`;
