queries = '''
    # A post object.
    post(id: String!): Post!
'''

mutations = '''
    # Publish a new post.
    pubPost(post: PostInput!): Post!
'''

types = '''
input PostInput {
    threadID: String!
    anonymous: Boolean!
    content: String!
    # Set referring PostIDs.
    refers: [String!]
}

type Post {
    id: String!
    anonymous: Boolean!
    author: String!
    content: String!
    createTime: Time!
    refers: [Post!]
    countOfRefered: Int!
}

type PostSlice {
    posts: [Post]!
    sliceInfo: SliceInfo!
}
'''
