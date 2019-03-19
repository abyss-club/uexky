import mongo from '~/utils/mongo';

// const index = (model, spec, opt) => model.collection.createIndex(
//   spec, { ...opt, background: true },
// );
function createIndexes() {
  const authCol = mongo.collection('auth');
  const notificationCol = mongo.collection('notification');
  const postCol = mongo.collection('post');
  const threadCol = mongo.collection('thread');
  const tokenCol = mongo.collection('token');
  const userCol = mongo.collection('user');
  const userAidCol = mongo.collection('userAid');
  const userPostsCol = mongo.collection('userPosts');

  return Promise.all([
    // index(AuthModel, { email: 1 }, { unique: true }),
    // index(AuthModel, { authCode: 1 }, {}),
    // index(AuthModel, { createdAt: 1 }, { expireAfterSeconds: 1200 }), // 20 minutes
    authCol.createIndexes([
      { key: { email: 1 }, unique: true },
      { key: { authCode: 1 } },
      { key: { createdAt: 1 }, expireAfterSeconds: 1200 },
    ]),
    // index(NotificationModel, { send_to: 1, type: 1, eventTime: 1 }, {}),
    // index(NotificationModel, { send_to_group: 1, type: 1, eventTime: 1 }, {
    //   partialFilterExpression: { send_to_group: { $exists: true } },
    // }),
    notificationCol.createIndexes([
      { key: { send_to: 1, type: 1, eventTime: 1 } },
      {
        key: { send_to_group: 1, type: 1, eventTime: 1 },
        partialFilterExpression: { send_to_group: { $exists: true } },
      },
    ]),

    // index(PostModel, { suid: 1 }, { unique: true }),
    // index(PostModel, { quoteSuids: 1 }, {}),

    postCol.createIndexes([
      { key: { suid: 1 }, unique: true },
      { key: { quoteSuids: 1 } },
    ]),

    // index(ThreadModel, { suid: 1 }, { unique: true }),
    // index(ThreadModel, { tags: 1, suid: -1 }, {}),

    threadCol.createIndexes([
      { key: { suid: 1 }, unique: true },
      { key: { tags: 1, suid: -1 } },
    ]),

    // index(TokenModel, { email: 1 }, { unique: true }),
    // index(TokenModel, { authToken: 1 }, { unique: true }),
    // index(TokenModel, { createdAt: 1 }, { expireAfterSeconds: 172800 }), // 20 days

    tokenCol.createIndexes([
      { key: { email: 1 }, unique: true },
      { key: { authToken: 1 }, unique: true },
      { key: { createdAt: 1 }, expireAfterSeconds: 172800 },
    ]),

    // index(UserModel, { email: 1 }, { unique: true }),

    userCol.createIndexes([
      { key: { email: 1 }, unique: true },
      { key: { name: 1 }, unique: true, partialFilterExpression: { name: { $type: 'string' } } },
    ]),

    // index(UserAidModel, { userId: 1, threadSuid: 1 }, { unique: true }),

    userAidCol.createIndexes([{ key: { userId: 1, threadSuid: 1 }, unique: true }]),

    // index(UserPostsModel, { userId: 1, threadSuid: 1, updatedAt: -1 }, {}),

    userPostsCol.createIndexes([{ key: { userId: 1, threadSuid: 1, updatedAt: -1 } }]),
  ]);
}

export default createIndexes;
