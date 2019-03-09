import mongoose from 'mongoose';
import sleep from 'sleep-promise';
import MongoMemoryServer, { MongoMemoryReplSet } from 'mongodb-memory-server';

import log from '~/utils/log';
// import createIndexes from '~/models/indexes';

let mongoServer;
let replSet;
const opts = {
  useNewUrlParser: true,
  useCreateIndex: true,
  useFindAndModify: false,
  autoIndex: false,
};

const startMongo = async () => {
  mongoServer = new MongoMemoryServer({ instance: { storageEngine: 'wiredTiger' } });
  const mongoUri = await mongoServer.getConnectionString();
  await mongoose.connect(mongoUri, opts, (err) => {
    if (err) log.error(err);
  });
  // await createIndexes();
  return mongoServer;
};

const startRepl = async () => {
  replSet = new MongoMemoryReplSet({
    instanceOpts: [
      { storageEngine: 'wiredTiger' },
    ],
  });
  await replSet.waitUntilRunning();
  const mongoUri = `${await replSet.getConnectionString()}?replicaSet=testset`;
  log.info(`connect test mongodb: ${mongoUri}`);

  await sleep(2000);

  await mongoose.connect(mongoUri, opts, (err) => {
    if (err) log.error(err);
  });
  // await createIndexes();
  return replSet;
};

// const stopMongo = () => {
//   mongoose.disconnect();
//   mongoServer.stop();
// };

// const stopRepl = () => {
//   mongoose.disconnect();
//   replSet.stop();
// };

export {
  startMongo, startRepl,
};
