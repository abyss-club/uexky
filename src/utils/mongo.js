import { MongoClient } from 'mongodb';
import { ParamsError } from '~/utils/error';

import log from '~/utils/log';

// eslint-disable-next-line import/no-mutable-exports
let db;
let mongoClient;

const options = {
  useNewUrlParser: true,
};

const connectDb = async (mongoUri, dbName) => {
  if (!mongoUri || !dbName) throw new ParamsError('Invalid mongoUri/dbName');
  mongoClient = new MongoClient(mongoUri, options);
  try {
    await mongoClient.connect();
  } catch (err) {
    log.error(err.stack);
  }
  db = mongoClient.db(dbName);
  return mongoClient;
};

const collection = (name) => {
  if (!db) {
    throw new ParamsError('MongoClient must be connected first.');
  }
  return db.collection(name);
};

const startSession = () => mongoClient.startSession();

export default {
  collection,
  startSession,
};

export {
  connectDb, db,
};
