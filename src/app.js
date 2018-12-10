import { ApolloServer, gql } from 'apollo-server-koa';
import * as Koa from 'koa';
import * as Router from 'koa-router';

import * as cors from '@koa/cors';

import { schema } from './schema';

const server = new ApolloServer({
  schema,
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

app.use(router.routes());
app.use(cors({ allowMethods: ['GET', 'OPTION', 'POST'] }));

export default app;
