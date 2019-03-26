import mongo from '~/utils/mongo';
import { expireTime } from './auth';

const index = async (collectionName, indexes) => {
  const idxes = indexes.map(idx => ({ ...idx, background: true }));
  await mongo.collection(collectionName).createIndexes(idxes);
};

function createIndexes() {
  return Promise.all([
    index('authCode', [
      { key: { email: 1 }, unique: true },
      { key: { authCode: 1 } },
      { key: { createdAt: 1 }, expireAfterSeconds: expireTime.code },
    ]),
    index('notification', [
      { key: { send_to: 1, type: 1, eventTime: 1 } },
      {
        key: { send_to_group: 1, type: 1, eventTime: 1 },
        partialFilterExpression: { send_to_group: { $exists: true } },
      },
    ]),
    index('post', [
      { key: { suid: 1 }, unique: true },
      { key: { quoteSuids: 1 } },
    ]),
    index('thread', [
      { key: { suid: 1 }, unique: true },
      { key: { tags: 1, suid: -1 } },
    ]),
    index('token', [
      { key: { email: 1 }, unique: true },
      { key: { authToken: 1 }, unique: true },
      { key: { createdAt: 1 }, expireAfterSeconds: expireTime.token },
    ]),
    index('user', [
      { key: { email: 1 }, unique: true },
      {
        key: { name: 1 },
        unique: true,
        partialFilterExpression: { name: { $type: 'string' } },
      },
    ]),
    index('userAid', [
      { key: { userId: 1, threadSuid: 1 }, unique: true },
    ]),
    index('userPosts', [
      { key: { userId: 1, threadSuid: 1, updatedAt: -1 } },
    ]),
  ]);
}

export default createIndexes;
export { index };
