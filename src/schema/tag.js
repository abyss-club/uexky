export default`
  extend type Query {
    # Containing mainTags and tagTree.
    tags: Tags!
  }

  type Tags {
    # Main tags are predefined manually.
    mainTags: [String!]!
    tree(
      "Keyword for filtering tags by tag name."
      query: String,
      "Maxmimum amount of returning \`MainTag\`s."
      limit: Int,
    ): [TagTreeNode!]
  }

  type TagTreeNode {
    mainTag: String!
    subTags: [String!]
  }
`;
