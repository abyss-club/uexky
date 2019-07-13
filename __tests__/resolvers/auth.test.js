import gql from 'graphql-tag';

import getRedis from '~/utils/redis';
import { mockUser, mutate, guestMutate } from '../__utils__/apolloClient';
import mockMailgun from '../__utils__/mailgun';
import startPg, { migrate } from '../__utils__/pgServer';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
  const redis = getRedis();
  await redis.flushall();
  redis.disconnect();
});

const mockEmail = mockUser.email;

const AUTH = gql`
  mutation Auth($email: String!) {
    auth(email: $email)
  }
`;

describe('Testing auth', () => {
  it('without context', async () => {
    mockMailgun();
    const { data } = await guestMutate({ mutation: AUTH, variables: { email: mockEmail } });
    expect(data.auth).toEqual(true);
  });
  it('with context', async () => {
    const { errors } = await mutate({ mutation: AUTH, variables: { email: mockEmail } });
    expect(errors[0].message).toEqual('Already signed in.');
  });
});
