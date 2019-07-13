import { Pool } from 'pg';
import { InternalError, ParamsError } from '~/utils/error';

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

const customErrors = new Set([
  'AuthError', 'InternalError', 'NotFoundError',
  'ParamsError', 'PermissionError',
]);
const handleError = (e) => {
  const name = e.name || '';
  if (customErrors.has(name)) {
    throw e;
  } else {
    log.error('transcation error', e.stack);
    throw new InternalError(`pg transcation error: ${e.stack}`);
  }
};

const query = async (text, params, client) => {
  let result;
  try {
    if (client) {
      result = await client.query(text, params);
    } else {
      result = await pgPool.query(text, params);
    }
  } catch (e) {
    throw new InternalError(`pg query '${text}' error: ${e.stack}`);
  }
  return result || {};
};

const doTransaction = async (transcation) => {
  const client = await pgPool.connect();
  try {
    await client.query('BEGIN');
    await transcation(client);
    await client.query('COMMIT');
  } catch (e) {
    await client.query('ROLLBACK');
    handleError(e);
  } finally {
    client.release();
  }
};

export {
  query, connectDb, doTransaction,
};
