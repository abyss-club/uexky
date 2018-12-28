import 'module-alias/register';
import { ApolloServer } from 'apollo-server-koa';
import Koa from 'koa';
import Router from 'koa-router';
import cors from '@koa/cors';

import mailguntest from './mailgun';
import { schema } from './schema';
import { genRandomStr } from './utils/uuid';
import { getUserByEmail } from './models/user';
import { addToAuth, getEmailByCode } from './models/auth';
import { genNewToken, getEmailByToken } from './models/token';

const server = new ApolloServer({
  schema,
  context: ({ ctx }) => {
    return { user: ctx.user };
  },
});

function authMiddleware() {
  return async (ctx, next) => {
    const token = ctx.cookies.get('token');
    if (ctx.url === '/graphql') {
      try {
        const email = await getEmailByToken(token);
        const user = await getUserByEmail(email);
        app.context.user = user;
      } catch (e) {
        if (e.authFail) app.context.user = null;
        else throw new Error(e);
      }
    }
    await next();
  };
}

const router = new Router();
// router.use('/graphql', authMiddleware());
router.get('/auth', async (ctx, next) => {
  if (!ctx.query.code || ctx.query.code.length !== 36) {
    ctx.body = 'Incorrect authentication code';
  } else {
    try {
      const email = await getEmailByCode(ctx.query.code);
      const token = await genNewToken(email);
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
app.use(authMiddleware());
app.use(cors({ allowMethods: ['GET', 'OPTION', 'POST'] }));
server.applyMiddleware({ app });
// mailguntest();

export default app;
