import request from 'supertest';
import app from '../src/app';

test('Hello API Request', async () => {
  const response = await request(app.callback()).get('/');
  expect(response.status).toEqual(200);
  expect(response.text).toEqual('Hello World!');
});
