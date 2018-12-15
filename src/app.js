import 'module-alias/register';
import { ApolloServer, gql } from 'apollo-server-koa';
import Koa from 'koa';
import Router from 'koa-router';
import cors from '@koa/cors';

import mailguntest from './mailgun';
import { schema } from './schema';
import { genRandomStr } from './utils/uuid';
import { getUserByEmail } from './models/user';
import { addToAuth, getEmailByCode } from './models/auth';
import { getTokenByEmail } from './models/token';

const server = new ApolloServer({
  schema,
  context: ({ req }) => {
    // get the user token from the headers
    const token = req.headers.authorization || '';
    console.log(req);
    // try to retrieve a user with the token
    const user = getUser(token);
    // add the user to the context
    return { user };
  },
});

const app = new Koa();
server.applyMiddleware({ app });

app.use(async (ctx, next) => {
  // Log the request to the console
  console.log('Url:', ctx.url);
  // Pass the request to the next middleware function
  await next();
});

const router = new Router();
router.get('/', async (ctx) => {
  ctx.body = 'Hello World!';
});
router.get('/auth', async (ctx, next) => {
  try {
    console.log(ctx.query.code);
    const email = await getEmailByCode(ctx.query.code);
    const token = await getTokenByEmail(email);
    const expiry = new Date(token.createdAt);
    expiry.setDate(expiry.getDate() + 20);
    console.log(expiry);
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
  await next();
});

app.use(router.routes());
app.use(cors({ allowMethods: ['GET', 'OPTION', 'POST'] }));
// mailguntest();

export default app;
