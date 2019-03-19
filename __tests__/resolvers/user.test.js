import gql from 'graphql-tag';

import startRepl from '../__utils__/mongoServer';
import { mockUser, query, mutate } from '../__utils__/apolloClient';

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

const SYNC_TAGS = gql`
  mutation SyncTags($tags: [String]!) {
    syncTags(tags: $tags) { tags }
  }
`;

const ADD_TAGS = gql`
  mutation AddTags($tags: [String!]!) {
    addSubbedTags(tags: $tags) { tags }
  }
`;

const DEL_TAGS = gql`
  mutation DelTags($tags: [String!]!) {
    delSubbedTags(tags: $tags) { tags }
  }
`;

const SET_NAME = gql`
  mutation SetName($name: String!) {
    setName(name: $name) { email, name }
  }
`;

const PROFILE = gql`
  query Profile {
    profile { email, name }
  }
`;

const subbedTags = ['MainA', 'SubA', 'SubB', 'MainB'];
const newSubbedTags = [...subbedTags, 'MainC'];
const dupSubbedTags = [...subbedTags, 'MainC', 'MainD'];

describe('Testing profile', () => {
  it('query profile', async () => {
    const { data } = await query({ query: PROFILE });
    expect(data.profile.email).toEqual(mockUser.email);
  });
});
describe('Testing modifying tags subscription', () => {
  it('sync tags', async () => {
    const { data } = await mutate({ mutation: SYNC_TAGS, variables: { tags: subbedTags } });
    expect(data.syncTags.tags).toEqual(subbedTags);
  });
});
describe('Testing adding tags subscription', () => {
  it('add tags', async () => {
    const { data } = await mutate({ mutation: ADD_TAGS, variables: { tags: ['MainC'] } });
    expect(JSON.stringify(data.addSubbedTags.tags)).toEqual(JSON.stringify(newSubbedTags));
  });
  it('add duplicated tags', async () => {
    const { data } = await mutate({ mutation: ADD_TAGS, variables: { tags: ['MainC', 'MainD'] } });
    expect(JSON.stringify(data.addSubbedTags.tags)).toEqual(JSON.stringify(dupSubbedTags));
  });
});
describe('Testing deleting tags subscription', () => {
  it('del tags', async () => {
    const { data } = await mutate({ mutation: DEL_TAGS, variables: { tags: ['MainD'] } });
    expect(JSON.stringify(data.delSubbedTags.tags)).toEqual(JSON.stringify(newSubbedTags));
  });
  it('del non-existing tags', async () => {
    const { data } = await mutate({ mutation: DEL_TAGS, variables: { tags: ['MainC', 'MainD'] } });
    expect(JSON.stringify(data.delSubbedTags.tags)).toEqual(JSON.stringify(subbedTags));
  });
});
describe('Testing setting name', () => {
  it('set name', async () => {
    const { data } = await mutate({ mutation: SET_NAME, variables: { name: mockUser.name } });
    expect(data.setName.email).toEqual(mockUser.email);
    expect(data.setName.name).toEqual(mockUser.name);
  });
  it('set name again', async () => {
    const { errors } = await mutate({ mutation: SET_NAME, variables: { name: mockUser.name } });
    expect(errors[0].message).toEqual('Name can only be set once.');
  });
});
