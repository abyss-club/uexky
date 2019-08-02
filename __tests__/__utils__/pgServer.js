import path from 'path';
import pgMigrate from 'node-pg-migrate';
import { connectDb } from '~/utils/pg';
import env from '~/utils/env';

const startPg = async () => {
  const pgPool = await connectDb(env.PG_URI);
  return pgPool;
};

const migrate = () => pgMigrate({
  databaseUrl: env.PG_URI,
  direction: 'up',
  dir: path.join(__dirname, '../../migrations'),
  migrationsTable: 'migrations',
  log: () => {}, // do not to log migration
}).catch((error) => {
  console.error(error);
});

export { migrate };
export default startPg;
