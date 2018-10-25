queries = '''
    # query unread notification count
    unreadNotiCount(): UnreadNotiCount!
    # query notification for user
    notification(type: String!, query: SliceQuery!): NotiSlice!
'''

types = '''
type UnreadNotiCount {
    system: Int!
    replied: Int!
    referred: Int!
}

type NotiSlice {
    system: [SystemNoti!]
    replied: [RepliedNoti!]
    referred: [ReferredNoti!]
    sliceInfo: SliceInfo!
}

type SystemNoti {
    id: String!
    type: String!
    eventTime: Time!
    hasRead: Boolean!
    title: String!
    content: String!
}

type RepliedNoti {
    id: String!
    type: String!
    eventTime: Time!
    hasRead: Boolean!
    thread: Thread!
    repliers: [String!]!
}

type ReferredNoti {
    id: String!
    type: String!
    eventTime: Time!
    hasRead: Boolean!
    thread: Thread!
    post: Post!
    referers: [String!]!
}
'''
