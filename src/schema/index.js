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
import UserResolver from '../resolvers/user';
import TagResolver from '../resolvers/tag';

const Query = `
  type Query {
    _empty: String
  }
`;

const resolvers = {};

export const schema = makeExecutableSchema({
  resolvers: Object.assign(resolvers, UserResolver, TagResolver),
  typeDefs: [Query, Base, Notification, Post, Tag, Thread, User],
});
