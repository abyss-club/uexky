import { Pool } from 'pg';
import { ParamsError } from '~/utils/error';

import log from '~/utils/log';

let pgPool;

const connectDb = async (pgUri) => {
  if (!pgUri) throw new ParamsError('Invalid mongoUri/dbName');
  pgPool = new Pool({ connectionString: pgUri });
  pgPool.on('error', (err) => {
    log.error('Unexpected error on idle client', err.stack);
    process.exit(-1);
  });
  return pgPool;
};

const query = (text, params) => pgPool.query(text, params);

const doTransaction = async (transcation) => {
  const client = pgPool.connect();
  try {
    await client.query('BEGIN');
    await transcation(client);
    await client.query('COMMIT');
  } catch (e) {
    await client.query('ROLLBACK');
    throw e;
  } finally {
    client.release();
  }
};

export {
  query, connectDb, doTransaction,
};
