import request from 'supertest';
import app, { endpoints, authMiddleware, configMiddleware } from '~/app';
// import { deflateSync } from 'zlib';
import Cookies from 'cookies';

import TokenModel from '~/models/token';
import UserModel from '~/models/user';
import ConfigModel from '~/models/config';
import startRepl from './__utils__/mongoServer';

jest.setTimeout(60000);

let replSet;
let mongoClient;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

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
  it('test not signed in', async () => {
    const middleware = authMiddleware();
    const ctx = { url: endpoints.graphql };
    ctx.cookies = new Cookies();
    await middleware(ctx, () => {});
    expect(ctx.user).toBeNull();
  });
  it('test signed in', async () => {
    // mock data
    const ctx = { url: endpoints.graphql };
    const mockEmail = 'test@example.com';
    const user = await UserModel().getUserByEmail(mockEmail);
    const token = await TokenModel(ctx).genNewToken(mockEmail);

    ctx.cookies = new Cookies();
    ctx.cookies.serialize('token', token);
    const middleware = authMiddleware();
    await middleware(ctx, () => {});

    expect(ctx.user.email).toEqual(user.email);
    expect(ctx.user).toEqual(user);
  });
});

describe('test config middleware', () => {
  it('get config', async () => {
    const middleware = configMiddleware();
    const ctx = { url: endpoints.graphql };
    await middleware(ctx, () => {});
    const config = await ConfigModel().getConfig();
    expect(ctx.config).toEqual(config);
  });
});
