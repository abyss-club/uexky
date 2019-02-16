import mongoose from 'mongoose';
import sleep from 'sleep-promise';
import MongoMemoryServer, { MongoMemoryReplSet } from 'mongodb-memory-server';

let mongoServer;
let replSet;
const opts = { useNewUrlParser: true };

const startMongo = async () => {
  mongoServer = new MongoMemoryServer({ instance: { storageEngine: 'wiredTiger' } });
  const mongoUri = await mongoServer.getConnectionString();
  await mongoose.connect(mongoUri, opts, (err) => {
    if (err) console.error(err);
  });
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

  await sleep(2000);

  await mongoose.connect(mongoUri, opts, (err) => {
    if (err) console.error(err);
  });

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
