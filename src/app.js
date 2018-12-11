import 'module-alias/register';
import { ApolloServer, gql } from 'apollo-server-koa';
import Koa from 'koa';
import Router from 'koa-router';
import cors from '@koa/cors';

import mailguntest from './mailgun';
import { schema } from './schema';

const server = new ApolloServer({
  schema,
  context: ({ req }) => {
    // get the user token from the headers
    const token = req.headers.authorization || '';
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
router.get('/auth', async (ctx) => {
  ctx.body = 'Authentication!';
});

app.use(router.routes());
app.use(cors({ allowMethods: ['GET', 'OPTION', 'POST'] }));
// mailguntest();

export default app;
