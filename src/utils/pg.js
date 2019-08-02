import { Pool } from 'pg';
import { InternalError, ParamsError } from '~/utils/error';

import log from '~/utils/log';

let pgPool;

const connectDb = async (pgUri) => {
  if (!pgUri) throw new ParamsError('Invalid mongoUri/dbName');
  pgPool = new Pool({ connectionString: pgUri });
  pgPool.on('error', (error) => {
    log.error('Unexpected error on idle client', error.stack);
    throw new InternalError(`Postgres client error: ${error.message}`);
  });
  return pgPool;
};

const customErrors = new Set([
  'AuthError', 'InternalError', 'NotFoundError',
  'ParamsError', 'PermissionError',
]);
const handleError = (error) => {
  const name = error.name || '';
  if (customErrors.has(name)) {
    throw error;
  } else {
    log.error(`transcation error: ${error.stack}`);
    throw new InternalError(`pg transcation error: ${error.message}`);
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
  } catch (error) {
    log.error(`pg query '${text}' with params ${params}, error: ${error.stack}`);
    throw new InternalError(`pg query error: ${error.message}`);
  }
  return result || {};
};

const doTransaction = async (transcation) => {
  const client = await pgPool.connect();
  try {
    await client.query('BEGIN');
    await transcation(client);
    await client.query('COMMIT');
  } catch (error) {
    await client.query('ROLLBACK');
    handleError(error);
  } finally {
    client.release();
  }
};

export {
  query, connectDb, doTransaction,
};
