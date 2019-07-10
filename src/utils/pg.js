import { Pool } from 'pg';
import { ParamsError, InternalError } from '~/utils/error';

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

const query = async (text, params) => {
  try {
    const result = await pgPool.query(text, params);
    return result;
  } catch (e) {
    throw new InternalError(`pg query '${text}' error: ${e.stack}`);
  }
};

const doTransaction = async (transcation) => {
  const client = await pgPool.connect();
  try {
    await client.query('BEGIN');
    await transcation(client);
    await client.query('COMMIT');
  } catch (e) {
    await client.query('ROLLBACK');
    log.error('transcation error', e.stack);
    throw new InternalError(`pg transcation error: ${e.stack}`);
  } finally {
    client.release();
  }
};

export {
  query, connectDb, doTransaction,
};
