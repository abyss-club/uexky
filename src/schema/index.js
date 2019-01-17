import { makeExecutableSchema } from 'graphql-tools';

import resolvers from '~/resolvers';

import base from './base';
import config from './config';
import notification from './notification';
import post from './post';
import tag from './tag';
import thread from './thread';
import user from './user';

const schema = makeExecutableSchema({
  resolvers,
  typeDefs: [base, config, notification, post, tag, thread, user],
});

export default schema;
