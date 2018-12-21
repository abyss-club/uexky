const tagQueries = `
  # Containing mainTags and tagTree.
  tags: Tags!
`;

const typeDef = `
  extend type Query {
    # Containing mainTags and tagTree.
    tags: Tags!
  }

  type Tags {
    # Main tags are predefined manually.
    mainTags: [String!]!
    tree(query: String): [TagTreeNode!]
  }

  type TagTreeNode {
    mainTag: String!
    subTags: [String!]
  }
`;

export { typeDef };
