import { gql } from 'apollo-server-koa';
import { makeExecutableSchema, addMockFunctionsToSchema } from 'graphql-tools';

import { typeDef as Base } from './base';
import { typeDef as Notification } from './notification';
import { typeDef as Post } from './post';
import { typeDef as Tag } from './tag';
import { typeDef as Thread } from './thread';
import { typeDef as User } from './user';

// import UserResolver, { profile } from '../resolvers/user';
// import TagResolver, { tags } from '../resolvers/tag';

import resolvers from '../resolvers';

// const Query = `
//   type Query {
//   }
// `;

export const schema = makeExecutableSchema({
  resolvers,
  typeDefs: [Base, Notification, Post, Tag, Thread, User],
});
