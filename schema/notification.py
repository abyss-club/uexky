queries = '''
    # query unread notification count
    unreadNotifCount(type: String): Int!
    # query notification for user
    notification(type: String!, query: SliceQuery!): NotifSlice!
'''

types = '''
type NotifSlice {
    notif: [Notif]!
    sliceInfo: SliceInfo!
}

type Notif {
    type String!
    releaseTime Time!
    hasRead Boolean!
    content String!
}
'''
