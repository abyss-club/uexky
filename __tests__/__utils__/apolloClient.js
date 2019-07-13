import { createTestClient } from 'apollo-server-testing';
import { ApolloServer } from 'apollo-server-koa';
import ConfigModel from '~/models/config';
import { createIdleRateLimiter } from '~/utils/rateLimit';

import UserModel from '~/models/user';
import schema from '~/schema';
import mockContext from './context';

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
    const auth = await UserModel.authContext({ email: mockUser.email });
    return { auth, config: await ConfigModel.getConfig(), limiter: createIdleRateLimiter() };
  },
});

const serverBeforeAuth = new ApolloServer({
  schema,
  context: async () => {
    const { auth } = await mockContext({});
    return { auth, config: await ConfigModel.getConfig(), limiter: createIdleRateLimiter() };
  },
});

const serverAltUser = new ApolloServer({
  schema,
  context: async () => {
    const { auth } = await mockContext({ email: mockAltUser.email });
    return { auth, config: await ConfigModel.getConfig(), limiter: createIdleRateLimiter() };
  },
});

const customClient = ({ email, name, role }) => {
  const cs = new ApolloServer({
    schema,
    context: async () => {
      const { auth } = await mockContext({ email, name, role });
      const config = await ConfigModel.getConfig();
      return { auth, config, limiter: createIdleRateLimiter() };
    },
  });
  const { query, mutate } = createTestClient(cs);
  return { query, mutate };
};

const { query, mutate } = createTestClient(server);
const { query: guestQuery, mutate: guestMutate } = createTestClient(serverBeforeAuth);
const { query: altQuery, mutate: altMutate } = createTestClient(serverAltUser);

export {
  mockUser, mockAltUser, query, mutate, guestQuery, guestMutate, altQuery, altMutate, customClient,
};
