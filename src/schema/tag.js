export default `
  extend type Query {
    # Main Tags.
    mainTags: [String!]!
    # Tags that are recommended.
    recommended: [String!]!
    # Searching tags by keyword.
    tags(
      # Search keyword.
      query: String,
      # Amount of tags returned.
      limit: Int,
    ): [Tag]!
  }

  type Tag {
    # Name of tag.
    name: String!
    # The tag is a Main Tag if true, Sub Tag otherwise.
    isMain: Boolean!
    # The tag which this tag belongs to.
    belongsTo: [String!]!
  }
`;
