import { ApolloServer, gql } from 'apollo-server-koa';
import { makeExecutableSchema } from 'graphql-tools';
import * as Koa from 'koa';
import * as Router from 'koa-router';
import * as mongoose from 'mongoose';

import * as cors from '@koa/cors';

import { schema } from './schema';

mongoose.connect('mongodb://localhost:27017/testing', { useNewUrlParser: true });

// Construct a schema, using GraphQL schema language
const typeDefs = gql`
  type Query {
    hello: String
  }
`;

// console.log(schema);

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

app.listen({ port: 5000 }, () =>
  console.log(`🚀 Server ready at http://localhost:5000${server.graphqlPath}`),
);
