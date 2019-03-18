import { MongoClient } from 'mongodb';
import { ParamsError } from '~/utils/error';

import log from '~/utils/log';
// import createIndexes from '~/models/indexes';

// let mongoServer;
// let replSet;
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
    console.error(err.stack);
  }
  db = mongoClient.db(dbName);
  return mongoClient;
};

const collection = (name) => {
  // console.log({ db });
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
