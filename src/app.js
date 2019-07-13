import { ApolloServer } from 'apollo-server-koa';
import Koa from 'koa';
import Router from 'koa-router';
import cors from '@koa/cors';

import schema from '~/schema';
import log from '~/utils/log';
import { authHandler, authMiddleware } from '~/auth';
import ConfigModel from '~/models/config';
import createRateLimiter, { createIdleRateLimiter } from '~/utils/rateLimit';

const endpoints = {
  graphql: '/graphql',
  auth: '/auth',
};

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
router.get(endpoints.auth, authHandler());
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
app.use(authMiddleware(endpoints.graphql));
app.use(cors({ allowMethods: ['GET', 'OPTION', 'POST'] }));
server.applyMiddleware({ app });

export default app;
export { endpoints, authMiddleware, configMiddleware };
