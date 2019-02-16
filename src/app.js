import { ApolloServer } from 'apollo-server-koa';
import Koa from 'koa';
import Router from 'koa-router';
import cors from '@koa/cors';

import schema from '~/schema';
import AuthModel from '~/models/auth';
import UserModel from '~/models/user';
import { config } from '~/models/config';
import createRateLimiter, { createIdleRateLimiter } from '~/utilities/rateLimit';
import TokenModel from '~/models/token';

const server = new ApolloServer({
  schema,
  context: ({ ctx }) => ({ user: ctx.user }),
});

function authMiddleware() {
  return async (ctx, next) => {
    const token = ctx.cookies.get('token');
    if (ctx.url === '/graphql') {
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
    if (ctx.url === '/graphql') {
      ctx.config = config();
    }
    await next();
  };
}

async function rateLimitMiddleware(ctx, next) {
  if (ctx.url === '/graphql') {
    const headerName = (await ctx.config.getRateLimit()).HTTPHeader;
    if (headerName !== '') {
      const ip = ctx.header.get(headerName);
      const { user } = ctx;
      if ((user) && (user.email) && (user.email !== '')) {
        ctx.limiter = createRateLimiter(ip, user.email);
      } else {
        ctx.limiter = createRateLimiter(ip);
      }
    } else {
      ctx.limiter = createIdleRateLimiter();
    }
  }
  await next();
}

const router = new Router();
router.get('/auth', async (ctx, next) => {
  if (!ctx.query.code || ctx.query.code.length !== 36) {
    ctx.body = 'Incorrect authentication code';
  } else {
    try {
      const email = await AuthModel.getEmailByCode(ctx.query.code);
      const token = await TokenModel.genNewToken(email);
      const expiry = new Date(token.createdAt);
      expiry.setDate(expiry.getDate() + 20);
      ctx.body = token;
      ctx.cookies.set('token', token.authToken, {
        path: '/',
        domain: process.env.API_DOMAIN,
        httpOnly: true,
        overwrite: true,
        expires: expiry,
      });
    } catch (e) {
      ctx.throw(401, e);
    }
  }
  await next();
});

const app = new Koa();
app.use(router.routes());
app.use(authMiddleware);
app.use(configMiddleware);
app.use(rateLimitMiddleware);
app.use(cors({ allowMethods: ['GET', 'OPTION', 'POST'] }));
server.applyMiddleware({ app });

export default app;
