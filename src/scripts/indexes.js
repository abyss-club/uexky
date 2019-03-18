import { connectDb } from '~/dbClient';
import env from '~/utils/env';
import log from '~/utils/log';
import createIndexes from '~/models/indexes';

// const result = (async () => {
//   const mongoClient = await connectDb(env.MONGODB_URI, env.MONGODB_DBNAME);
//   await createIndexes();
//   await mongoClient.close();
// })();

const initIndex = async () => {
  const mongoClient = await connectDb(env.MONGODB_URI, env.MONGODB_DBNAME);
  await createIndexes();
  await mongoClient.close();
};

try {
  initIndex();
} catch (e) {
  log.error(e);
}

// log.info(`${result}`);
