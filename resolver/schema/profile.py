queries = '''
    # A user profile object.
    profile(): User!
    # query notification for user
    notification(): NotificationSlice!
'''

mutations = '''
    # Register/Login via email address. An email containing login info will be sent to the provided email address.
    auth(email: String!): Boolean!
    # Set the Name of user.
    setName(name: String!): User!
    # Save/Add/Del tags subscribed by user.
    syncTags(tags: [String]!): User!
    addSubbedTags(tags: [String!]!): User!
    delSubbedTags(tags: [String!]!): User!
'''

types = '''
type User {
    email: String!
    # The Name of user. Required when not posting anonymously.
    name: String
    # Tags saved by user.
    tags: [String!]
}

type Notification {
    releaseTime Time!
    type String!
    hasRead Boolean!
    content String!
}

type NotificationSlice {
    unreadCount Int!
    notification [Notification]!
    sliceInfo SliceInfo!
}
'''
