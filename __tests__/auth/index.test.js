import Koa from 'koa';
import request from 'supertest';

import { authMiddleware, authHandler } from '~/auth';
import Token from '~/auth/token';
import { Base64 } from '~/uid';
import Code, { expireTime } from '~/auth/code';
import getRedis from '~/utils/redis';
import log from '~/utils/log';
import env from '~/utils/env';
import mockMailgun from '../__utils__/mailgun';

afterAll(async () => {
  const redis = getRedis();
  await redis.flushall();
  redis.disconnect();
});

it('test auth middleware/without token', async () => {
  let user;
  const app = new Koa();
  app.use(authMiddleware('/test'));
  app.use((ctx) => {
    user = ctx.user || {};
    ctx.body = user.email || '';
  });

  const response = await request(app.callback()).get('/test');
  expect(response.text).toEqual('');
});

it('test auth middleware/with token', async () => {
  let user;
  const app = new Koa();
  app.use(authMiddleware('/test'));
  app.use((ctx) => {
    user = ctx.user || {};
    ctx.body = user.email || '';
  });
  const mockEmail = 'example@example.com';
  const token = await Token.genNewToken(mockEmail);

  const response = await request(app.callback())
    .get('/test')
    .set('Cookie', `token=${token}`);
  expect(response.text).toEqual(mockEmail);
  expect(user.email).toEqual(mockEmail);
});

it('test auth handler/without code', async () => {
  const app = new Koa();
  app.use(authHandler());
  const response = await request(app.callback()).get('/auth');
  expect(response.status).toEqual(400);
});

it('test auth handler/wrong code', async () => {
  const app = new Koa();
  app.use(authHandler());
  const code = Base64.randomString(36);
  const response = await request(app.callback()).get(`/auth?code=${code}`);
  expect(response.status).toEqual(401);
});

it('test auth handler/with code', async () => {
  const app = new Koa();
  app.use(authHandler());
  const mockEmail = 'example@example.com';

  mockMailgun();
  const code = await Code.addToAuth(mockEmail);

  env.DOMAIN = 'uexky.com';
  env.PROTO = 'https';
  const response = await request(app.callback())
    .get(`/auth?code=${code}`)
    .expect(302)
    .expect('Cache-Control', 'no-cache, no-store')
    .expect('Location', `${env.PROTO}://${env.DOMAIN}`);
  const setCookie = response.headers['set-cookie'];
  const cookie = parseCookie(setCookie[0]);
  expect(cookie.Path).toEqual('/');
  expect(cookie['Max-Age']).toEqual(`${expireTime.token}`);
  expect(cookie.Domain).toEqual(env.DOMAIN);
  expect(cookie.HttpOnly).toBeTruthy();
  expect(cookie.Secure).toBeTruthy();
});

function parseCookie(cookiesHeader) {
  log.debug(`cookies header: ${cookiesHeader}`);
  const cookies = {};
  cookiesHeader.split(';').forEach((pair) => {
    const words = pair.split('=');
    if (words.length === 1) {
      cookies[words[0]] = true;
    } else if (words.length === 2) {
      const [key, value] = words;
      cookies[key] = value;
    }
  });
  log.debug('cookies', cookies);
  return cookies;
}
