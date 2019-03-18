import { createTestClient } from 'apollo-server-testing';
import { ApolloServer } from 'apollo-server-koa';
import ConfigModel from '~/models/config';
import { createIdleRateLimiter } from '~/utils/rateLimit';

import UserModel from '~/models/user';
import schema from '~/schema';

const mockUser = {
  email: 'test@example.com',
  name: 'testUser',
};

const mockAltUser = {
  email: 'alt@example.com',
  name: 'altUser',
};

const server = new ApolloServer({
  schema,
  context: async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    return { user, config: await ConfigModel().getConfig(), limiter: createIdleRateLimiter() };
  },
});

const serverBeforeAuth = new ApolloServer({
  schema,
  context: async () => ({
    user: null, config: await ConfigModel().getConfig(), limiter: createIdleRateLimiter(),
  }),
});

const serverAltUser = new ApolloServer({
  schema,
  context: async () => {
    const user = await UserModel().getUserByEmail(mockAltUser.email);
    return { user, config: await ConfigModel().getConfig(), limiter: createIdleRateLimiter() };
  },
});

const { query, mutate } = createTestClient(server);
const { query: guestQuery, mutate: guestMutate } = createTestClient(serverBeforeAuth);
const { query: altQuery, mutate: altMutate } = createTestClient(serverAltUser);

export {
  mockUser, mockAltUser, query, mutate, guestQuery, guestMutate, altQuery, altMutate,
};
