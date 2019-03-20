import request from 'supertest';
import app, { endpoints } from '~/app';

describe('Testing paths', () => {
  it('Get /invalid', async () => {
    const response = await request(app.callback()).get('/invalid');
    expect(response.status).toEqual(404);
    expect(response.text).toEqual('Not Found');
  });
});

describe('Testing auth', () => {
  it(`Plain request to ${endpoints.auth}`, async () => {
    const response = await (await request(app.callback())).get(endpoints.auth);
    expect(response.status).toEqual(400);
    expect(response.text).toEqual('验证信息格式错误');
  });
});
