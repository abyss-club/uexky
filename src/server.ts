import * as Koa from 'koa';
import * as Router from 'koa-router';
import { ApolloServer, gql } from 'apollo-server-koa';

const app = new Koa();

// Construct a schema, using GraphQL schema language
const typeDefs = gql`
  type Query {
    hello: String
  }
`;

// Provide resolver functions for your schema fields
const resolvers = {
  Query: {
    hello: () => 'Hello world!',
  },
};

const server = new ApolloServer({ typeDefs, resolvers });

server.applyMiddleware({ app });

app.use(async (ctx, next) => {
    // Log the request to the console
    console.log('Url:', ctx.url);
    // Pass the request to the next middleware function
    await next();
});

const router = new Router();
router.get('/*', async (ctx) => {
    ctx.body = 'Hello World!';
});
app.use(router.routes());

app.listen({ port: 5000 }, () =>
  console.log(`🚀 Server ready at http://localhost:5000${server.graphqlPath}`),
);
