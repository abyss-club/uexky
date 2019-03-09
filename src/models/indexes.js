import mongo from '~/utils/mongo';

const index = (collectionName, spec, opt) => {
  mongo.collection(collectionName).createIndex(spec, { ...opt, background: true });
};

function createIndexes() {
  return Promise.all([
    index('authCode', { email: 1 }, { unique: true }),
    index('authCode', { authCode: 1 }, {}),
    index('authCode', { createdAt: 1 }, { expireAfterSeconds: 1200 }), // 20 minutes

    index('notification', { send_to: 1, type: 1, eventTime: 1 }, {}),
    index('notification', { send_to_group: 1, type: 1, eventTime: 1 }, {
      partialFilterExpression: { send_to_group: { $exists: true } },
    }),

    index('post', { suid: 1 }, { unique: true }),
    index('post', { quoteSuids: 1 }, {}),

    index('thread', { suid: 1 }, { unique: true }),
    index('thread', { tags: 1, suid: -1 }, {}),


    index('token', { email: 1 }, { unique: true }),
    index('token', { authToken: 1 }, { unique: true }),
    index('token', { createdAt: 1 }, { expireAfterSeconds: 172800 }), // 20 days

    index('user', { email: 1 }, { unique: true }),
    index('user', { name: 1 }, {
      unique: true,
      partialFilterExpression: { name: { $type: 'string' } },
    }),

    index('userAid', { userId: 1, threadSuid: 1 }, { unique: true }),

    index('userPosts', { userId: 1, threadSuid: 1, updatedAt: -1 }, {}),
  ]);
}

export default createIndexes;
