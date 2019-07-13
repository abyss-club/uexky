import { connectDb } from '~/utils/pg';
import path from 'path';
import env from '~/utils/env';
import pgMigrate from 'node-pg-migrate';

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
}).catch((e) => {
  console.error(e);
});

export { migrate };
export default startPg;
