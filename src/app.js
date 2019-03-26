import { ApolloServer } from 'apollo-server-koa';
import Koa from 'koa';
import Router from 'koa-router';
import cors from '@koa/cors';

import schema from '~/schema';
import log from '~/utils/log';
import env from '~/utils/env';
import AuthModel, { expireTime } from '~/models/auth';
import UserModel from '~/models/user';
import ConfigModel from '~/models/config';
import createRateLimiter, { createIdleRateLimiter } from '~/utils/rateLimit';
import TokenModel from '~/models/token';

const endpoints = {
  graphql: '/graphql',
  auth: '/auth',
};

function setCookie(ctx, token) {
  const opts = {
    path: '/',
    domain: env.DOMAIN,
    maxAge: expireTime.token,
    httpOnly: true,
    overwrite: true,
  };
  if (env.PROTO === 'https') {
    opts.secure = true;
  }
  ctx.cookies.set('token', token, opts);
}

function authMiddleware() {
  return async (ctx, next) => {
    const token = ctx.cookies.get('token') || '';
    if ((ctx.url === endpoints.graphql) && (token !== '')) {
      try {
        const email = await TokenModel.getEmailByToken(token);
        const user = await UserModel.getUserByEmail(email, true);
        ctx.user = user;
        setCookie(ctx, token);
      } catch (e) {
        if (e.authError) ctx.user = null;
        else throw new Error(e);
      }
    }
    await next();
  };
}

function configMiddleware() {
  return async (ctx, next) => {
    if (ctx.url === endpoints.graphql) {
      ctx.config = await ConfigModel.getConfig();
    }
    await next();
  };
}

function rateLimitMiddleware() {
  return async (ctx, next) => {
    if (ctx.url === endpoints.graphql) {
      const headerName = ctx.config.rateLimit.httpHeader;
      if (headerName !== '') {
        const ip = ctx.header.get(headerName);
        const { user } = ctx;
        if ((user) && (user.email) && (user.email !== '')) {
          ctx.limiter = createRateLimiter(ctx.config, ip, user.email);
        } else {
          ctx.limiter = createRateLimiter(ctx.config, ip);
        }
      } else {
        ctx.limiter = createIdleRateLimiter();
      }
    }
    await next();
  };
}

function logMiddleware() {
  return async (ctx, next) => {
    const start = new Date();
    try {
      await next();
    } catch (e) {
      log.error(e);
      throw e;
    }
    const stop = new Date();
    log.info('request', {
      href: ctx.href,
      take: `${stop - start}ms`,
      method: ctx.method,
      status: ctx.status,
    });
  };
}

const router = new Router();
router.get(endpoints.auth, async (ctx, next) => {
  if (!ctx.query.code || ctx.query.code.length !== 36) {
    ctx.throw(400, '验证信息格式错误');
  } else {
    try {
      const email = await AuthModel.getEmailByCode(ctx.query.code);
      const token = await TokenModel.genNewToken(email);
      setCookie(ctx, token);
      ctx.response.header.set('Location', `${env.PROTO}://${env.DOMAIN}`);
      ctx.response.header.set('Cache-Control', 'no-cache, no-store');
      ctx.response.status = 302;
    } catch (e) {
      log.error(e);
      ctx.throw(401, '验证信息错误或已失效');
    }
  }
  await next();
});

const server = new ApolloServer({
  schema,
  context: ({ ctx }) => ({
    config: ctx.config,
    user: ctx.user,
    limiter: ctx.limiter,
  }),
});

const app = new Koa();

// use middlewares from top to bottom;
app.use(router.routes());
app.use(logMiddleware());
app.use(configMiddleware());
app.use(rateLimitMiddleware());
app.use(authMiddleware());
app.use(cors({ allowMethods: ['GET', 'OPTION', 'POST'] }));
server.applyMiddleware({ app });

export default app;
export { endpoints };
