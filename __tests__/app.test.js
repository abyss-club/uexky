import request from 'supertest';
import app, { endpoints, authMiddleware } from '~/app';
// import { deflateSync } from 'zlib';
import Cookies from 'cookies';

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

describe('test auth middleware', () => {
  it('test auth middleware', async () => {
    const middleware = authMiddleware();
    const ctx = {};
    ctx.cookies = new Cookies();
    const token = '';
    ctx.cookies.serialize('token', token);
    await middleware(ctx, () => {});
  });
});
