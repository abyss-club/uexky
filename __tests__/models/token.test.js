import TokenModel from '~/models/token';
import startPg, { migrate } from '../__utils__/pgServer';

let ctx;
let pgPool;
// let db;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

describe('Testing token', () => {
  const mockEmail = 'test@example.com';
  it('validate token by email', async () => {
    const tokenResult = await TokenModel(ctx).genNewToken(mockEmail);
    const emailResult = await TokenModel(ctx).getEmailByToken(tokenResult);
    expect(emailResult).toEqual(mockEmail);
  });
});
