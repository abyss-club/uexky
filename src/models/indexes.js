import AuthModel from '~/models/auth';
import NotificationModel from '~/models/notification';
import PostModel from '~/models/post';
import ThreadModel from '~/models/thread';
import TokenModel from '~/models/token';
import UserModel from '~/models/user';
import UserAidModel from '~/models/userAid';
import UserPostsModel from '~/models/userPosts';

const index = (model, spec, opt) => model.collection.createIndex(
  spec, { ...opt, background: true },
);

async function createIndexes() {
  await Promise.all([
    index(AuthModel, { email: 1 }, { unique: true }),
    index(AuthModel, { authCode: 1 }, {}),
    index(AuthModel, { createdAt: 1 }, { expireAfterSeconds: 1200 }),

    index(NotificationModel, { send_to: 1, type: 1, eventTime: 1 }, {}),
    index(NotificationModel, { send_to_group: 1, type: 1, eventTime: 1 }, {
      partialFilterExpression: { send_to_group: { $exist: true } },
    }),

    index(PostModel, { suid: 1 }, { unique: true }),
    index(PostModel, { quoteSuids: 1 }, {}),

    index(ThreadModel, { suid: 1 }, { unique: true }),
    index(ThreadModel, { tags: 1, suid: -1 }, {}),

    index(TokenModel, { email: 1 }, { unique: true }),
    index(TokenModel, { authToken: 1 }, { unique: true }),
    index(TokenModel, { createdAt: 1 }, { expireAfterSeconds: 172800 }), // 20 days

    index(UserModel, { email: 1 }, { unique: true }),

    index(UserAidModel, { userId: 1, threadSuid: 1 }, { unique: true }),

    index(UserPostsModel, { userId: 1, threadSuid: 1, updatedAt: -1 }, {}),
  ]);
}

export default createIndexes;
