import gql from 'graphql-tag';
import mongoose from 'mongoose';

import { startMongo } from '../__utils__/mongoServer';
import { mockUser, mutate, guestMutate } from '../__utils__/apolloClient';

let mongoServer;

beforeAll(async () => {
  mongoServer = await startMongo();
});

afterAll(() => {
  mongoose.disconnect();
  mongoServer.stop();
});

const mockEmail = mockUser.email;

const AUTH = gql`
  mutation Auth($email: String!) {
    auth(email: $email)
  }
`;

describe('Testing auth', () => {
  it('without context', async () => {
    const { data } = await guestMutate({ mutation: AUTH, variables: { email: mockEmail } });
    expect(data.auth).toEqual(true);
  });
  it('with context', async () => {
    const { errors } = await mutate({ mutation: AUTH, variables: { email: mockEmail } });
    expect(errors[0].message).toEqual('Already signed in.');
  });
});
