import request from 'supertest';
import app from '~/app';

// import startPg, { migrate } from './__utils__/pgServer';
//
// let pgPool;
//
// beforeAll(async () => {
//   await migrate();
//   pgPool = await startPg();
// });
//
// afterAll(async () => {
//   await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
//   pgPool.end();
// });

describe('Testing paths', () => {
  it('Get /invalid', async () => {
    const response = await request(app.callback()).get('/invalid');
    await expect(response.status).toEqual(404);
    await expect(response.text).toEqual('Not Found');
  });
});
