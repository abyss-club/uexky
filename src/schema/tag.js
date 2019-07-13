export default `
  extend type Query {
    # main tags
    mainTags: [String!]!
    # Containing mainTags and tagTree.
    tags(query: String, limit: Int): [Tag]!
  }

  type Tag {
    # tag's name
    name: String!
    isMain: Boolean!
    # belongsTo which main tag
    belongsTo: [String!]!
  }
`;
