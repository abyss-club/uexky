const tagQueries = `
  # Containing mainTags and tagTree.
  tags: Tags!
`;

const tagTypes = `
  type Tags {
    # Main tags are predefined manually.
    mainTags: [String!]!
    # Recommended tags are picked manually.
    recommended: [String!]!
    tree(query: String): [TagTreeNode!]
  }

  type TagTreeNode {
    mainTag: String!
    subTags: [String!]
  }
`;

export { tagQueries, tagTypes };
export default `${tagQueries}
${tagTypes}`;
