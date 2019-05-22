import { connectDb } from '~/utils/pg';
import path from 'path';
import env from '~/utils/env';
import pgMigrate from 'node-pg-migrate';

const startPg = async () => {
  const pgPool = await connectDb(env.PG_URI, env.PG_DBNAME);
  return pgPool;
};

const migrate = async () => pgMigrate({
  databaseUrl: `${env.PG_URI}/${env.PG_DBNAME}`,
  direction: 'up',
  dir: path.join(__dirname, '../../migrations'),
  migrationsTable: 'pgmigrations',
});

export { migrate };
export default startPg;
