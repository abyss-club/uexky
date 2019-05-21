import { Pool } from 'pg';
import { ParamsError } from '~/utils/error';

import log from '~/utils/log';

// eslint-disable-next-line import/no-mutable-exports
let db;
let pgPool;

const connectDb = async (pgUri, dbName) => {
  if (!pgUri) throw new ParamsError('Invalid mongoUri/dbName');
  pgPool = new Pool({ connectionString: `${pgUri}/${dbName}` });
  pgPool.on('error', (err) => {
    log.error('Unexpected error on idle client', err.stack);
    process.exit(-1);
  });
  return pgPool;
};

const query = (text, params) => pgPool.query(text, params);

const getClient = () => pgPool.connect();
// const collection = (name) => {
//   if (!db) {
//     throw new ParamsError('MongoClient must be connected first.');
//   }
//   return db.collection(name);
// };
//
// const startSession = () => mongoClient.startSession();

export default {
  query,
  getClient,
};

export {
  connectDb, db,
};
