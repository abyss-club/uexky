import { ApolloServer } from 'apollo-server-koa';
import Koa from 'koa';
import Router from 'koa-router';
import cors from '@koa/cors';

import schema from '~/schema';
import log from '~/utils/log';
import env from '~/utils/env';
import AuthModel from '~/models/auth';
import UserModel from '~/models/user';
import { getConfig } from '~/models/config';
import createRateLimiter, { createIdleRateLimiter } from '~/utils/rateLimit';
import TokenModel from '~/models/token';

const server = new ApolloServer({
  schema,
  context: ({ ctx }) => ({
    config: ctx.config,
    user: ctx.user,
    limiter: ctx.limiter,
  }),
});

const endpoints = {
  graphql: '/graphql',
  auth: '/auth',
};

function authMiddleware() {
  return async (ctx, next) => {
    const token = ctx.cookies.get('token');
    if (ctx.url === endpoints.graphql) {
      try {
        const email = await TokenModel.getEmailByToken(token);
        const user = await UserModel.getUserByEmail(email);
        ctx.user = user;
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
      ctx.getConfig = getConfig();
    }
    await next();
  };
}

function rateLimitMiddleware() {
  return async (ctx, next) => {
    if (ctx.url === endpoints.graphql) {
      const config = await ctx.getConfig();
      const headerName = config.rateLimit.httpHeader;
      if (headerName !== '') {
        const ip = ctx.header.get(headerName);
        const { user } = ctx;
        if ((user) && (user.email) && (user.email !== '')) {
          ctx.limiter = createRateLimiter(config, ip, user.email);
        } else {
          ctx.limiter = createRateLimiter(config, ip);
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
    ctx.body = 'Incorrect authentication code';
  } else {
    try {
      const email = await AuthModel.getEmailByCode(ctx.query.code);
      const token = await TokenModel.genNewToken(email);
      const expiry = new Date(token.createdAt);
      expiry.setDate(expiry.getDate() + 20); // TODO(tangwenhan): time setting
      ctx.body = token;
      ctx.cookies.set('token', token.authToken, {
        path: '/',
        domain: env.API_DOMAIN,
        httpOnly: true,
        overwrite: true,
        expires: expiry,
      });
      ctx.response.status = 302;
      ctx.response.header.set('Location', `${env.PROTO}://${env.DOMAIN}`);
      ctx.response.header.set('Cache-Control', 'no-cache, no-store');
    } catch (e) {
      log.error(e);
      ctx.throw(401, e);
    }
  }
  await next();
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
