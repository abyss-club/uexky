// import { MongoClient } from 'mongodb';
import { connectDb } from '~/dbClient';
import sleep from 'sleep-promise';
import { MongoMemoryReplSet } from 'mongodb-memory-server';

import log from '~/utils/log';
// import createIndexes from '~/models/indexes';

// let mongoServer;
// let replSet;

// const options = {
//   useNewUrlParser: true,
// };

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

// const startMongo = async () => {
//   const mongoServer = new MongoMemoryServer({ instance: { storageEngine: 'wiredTiger' } });
//   const mongoUri = await mongoServer.getConnectionString();
//   const mongoClient = new MongoClient(mongoUri, options);
//
//   try {
//     await mongoClient.connect();
//   } catch (err) {
//     console.error(err.stack);
//   }
//   const db = mongoClient.db(await mongoServer.getDbName());
//   return { mongoClient, mongoServer, db };
// };

// const startRepl = async () => {
//   const replSet = new MongoMemoryReplSet({
//     instanceOpts: [
//       { storageEngine: 'wiredTiger' },
//     ],
//   });
//   await replSet.waitUntilRunning();
//   const mongoUri = `${await replSet.getConnectionString()}?replicaSet=testset`;
//   log.info(`connect test mongodb: ${mongoUri}`);
//
//   await sleep(2000);
//
//   const mongoClient = new MongoClient(mongoUri, options);
//   console.log(typeof mongoClient);
//   try {
//     await mongoClient.connect();
//   } catch (err) {
//     console.log(err.stack);
//   }
//
//   const db = mongoClient.db(await replSet.getDbName());
//   return { replSet, mongoClient, db };
// };

// const stopMongo = () => {
//   mongoServer.stop();
// };

// const stopRepl = () => {
//   replSet.stop();
// };

export { startRepl };
