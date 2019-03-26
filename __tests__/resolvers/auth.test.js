import gql from 'graphql-tag';

import startRepl from '../__utils__/mongoServer';
import { mockUser, mutate, guestMutate } from '../__utils__/apolloClient';
import { mockMailgun } from '~/utils/authMail';

jest.setTimeout(60000);
let replSet;
let mongoClient;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

const mockEmail = mockUser.email;

const AUTH = gql`
  mutation Auth($email: String!) {
    auth(email: $email)
  }
`;

const mailgun = {
  mail: null,
  messages: function messages() {
    const that = this;
    return {
      send(mail, fallback) {
        that.mail = mail;
        fallback(null, 'success');
      },
    };
  },
};

describe('Testing auth', () => {
  it('without context', async () => {
    mockMailgun(mailgun);
    const { data } = await guestMutate({ mutation: AUTH, variables: { email: mockEmail } });
    expect(data.auth).toEqual(true);
  });
  it('with context', async () => {
    const { errors } = await mutate({ mutation: AUTH, variables: { email: mockEmail } });
    expect(errors[0].message).toEqual('Already signed in.');
  });
});
