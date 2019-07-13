import gql from 'graphql-tag';

import { query as pgq } from '~/utils/pg';
import { mockUser, query, mutate } from '../__utils__/apolloClient';
import startPg, { migrate } from '../__utils__/pgServer';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

const SYNC_TAGS = gql`
  mutation SyncTags($tags: [String]!) {
    syncTags(tags: $tags) { tags }
  }
`;

const ADD_TAG = gql`
  mutation AddTag($tag: String!) {
    addSubbedTag(tag: $tag) { tags }
  }
`;

const DEL_TAG = gql`
  mutation DelTags($tag: String!) {
    delSubbedTag(tag: $tag) { tags }
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

describe('Testing profile', () => {
  it('query profile', async () => {
    const { data } = await query({ query: PROFILE });
    expect(data.profile.email).toEqual(mockUser.email);
  });
});
describe('Testing modifying tags subscription', () => {
  it('parpare data', async () => {
    await pgq('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['MainA', true]);
    await pgq('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['MainB', true]);
    await pgq('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['MainC', true]);
    await pgq('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['MainD', true]);
    await pgq('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['SubA', false]);
    await pgq('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['SubB', false]);
  });
  it('sync tags', async () => {
    const { data, errors } = await mutate({ mutation: SYNC_TAGS, variables: { tags: subbedTags } });
    expect(errors).toBeUndefined();
    expect(data.syncTags.tags).toEqual(subbedTags);
  });
  it('add tags', async () => {
    const { data, errors } = await mutate({ mutation: ADD_TAG, variables: { tag: 'MainC' } });
    expect(errors).toBeUndefined();
    expect(JSON.stringify(data.addSubbedTag.tags)).toEqual(JSON.stringify(newSubbedTags));
  });
  it('del tags', async () => {
    const { data, errors } = await mutate({ mutation: DEL_TAG, variables: { tag: 'MainC' } });
    expect(errors).toBeUndefined();
    expect(JSON.stringify(data.delSubbedTag.tags)).toEqual(JSON.stringify(subbedTags));
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
