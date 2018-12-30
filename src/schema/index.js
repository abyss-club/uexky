import { makeExecutableSchema } from 'graphql-tools';

import resolvers from '~/resolvers';

import base from './base';
import notification from './notification';
import post from './post';
import tag from './tag';
import thread from './thread';
import user from './user';

const schema = makeExecutableSchema({
  resolvers,
  typeDefs: [base, notification, post, tag, thread, user],
});

export default schema;
