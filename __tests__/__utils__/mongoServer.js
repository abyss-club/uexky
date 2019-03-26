// import { MongoClient } from 'mongodb';
import { connectDb } from '~/utils/mongo';
import sleep from 'sleep-promise';
import { MongoMemoryReplSet } from 'mongodb-memory-server';

import log from '~/utils/log';

const startRepl = async () => {
  const replSet = new MongoMemoryReplSet({
    replSet: {
      count: 3,
      storageEngine: 'wiredTiger',
    },
  });
  await replSet.waitUntilRunning();
  await sleep(5000);
  const mongoUri = `${await replSet.getConnectionString()}?replicaSet=testset`;
  log.info(`connect test mongodb: ${mongoUri}`);

  // process.env.MONGODB_URI = mongoUri;
  // process.env.MONGODB_DBNAME = await replSet.getDbName();
  const dbName = await replSet.getDbName();

  const mongoClient = await connectDb(mongoUri, dbName);
  return { replSet, mongoClient };
};

export default startRepl;
