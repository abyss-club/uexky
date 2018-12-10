import { gql } from 'apollo-server-koa';
import { makeExecutableSchema } from 'graphql-tools';

import { typeDef as Base } from './base';
import { typeDef as Notification } from './notification';
import { typeDef as Post } from './post';
import { typeDef as Tag } from './tag';
import { typeDef as Thread } from './thread';
import { typeDef as User } from './user';

const Query = `
  type Query {
    _empty: String
  }
`;

const resolvers = {
  Query: {
    test: () => 'Hello world!',
  },
};

export const schema = makeExecutableSchema({
  resolvers,
  typeDefs: [Query, Base, Notification, Post, Tag, Thread, User],
});
