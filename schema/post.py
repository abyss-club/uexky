queries = '''
    # A post object.
    post(id: String!): Post!
'''

mutations = '''
    # Publish a new post.
    pubPost(post: PostInput!): Post!
'''

types = '''
# Input object describing a Post to be published.
input PostInput {
    threadID: String!
    anonymous: Boolean!
    content: String!
    # Set quoting PostIDs.
    quotes: [String!]
}

# Object describing a Post.
type Post {
    id: String!
    anonymous: Boolean!
    author: String!
    content: String!
    createTime: Time!
    quotes: [Post!]
    quoteCount: Int!
}

# PostSlice object is for selecting specific 'slice' of Post objects to return. Affects returned SliceInfo.
type PostSlice {
    posts: [Post]!
    sliceInfo: SliceInfo!
}
'''
