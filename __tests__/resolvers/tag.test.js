import gql from 'graphql-tag';

import ThreadModel from '~/models/thread';
import { query as pgq } from '~/utils/pg';
import startPg, { migrate } from '../__utils__/pgServer';
import mockContext from '../__utils__/context';
import { customClient } from '../__utils__/apolloClient';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

describe('Insert Tags', () => {
  const mockEmail = 'test@uexky.com';
  const { query } = customClient({ email: mockEmail });
  it('parpare data', async () => {
    await pgq('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['MainA', true]);
    const ctx = await mockContext({ email: mockEmail });
    await ThreadModel.new({
      ctx,
      thread: {
        anonymous: true,
        content: 'Test Content',
        mainTag: 'MainA',
        subTags: ['SubA', 'SubB'],
        title: 'TestTitle',
      },
    });
  });
  it('get main tags', async () => {
    const GET_MAIN_TAGS = gql`
      query {
        mainTags
      }
    `;
    const { data, errors } = await query({ query: GET_MAIN_TAGS });
    expect(errors).toBeUndefined();
    expect(data.mainTags).toEqual(['MainA']);
  });
  it('query tags', async () => {
    const QUERY_TAGS = gql`
      query Tags($query: String!) {
        tags(query: $query, limit: 10) {
          name, isMain, belongsTo
        }
      }
    `;
    const { data, errors } = await query({
      query: QUERY_TAGS,
      variables: { query: 'A' },
    });
    expect(errors).toBeUndefined();
    expect(data.tags.length).toEqual(2);
    expect(data.tags).toContainEqual({
      name: 'MainA', isMain: true, belongsTo: [],
    });
    expect(data.tags).toContainEqual({
      name: 'SubA', isMain: false, belongsTo: ['MainA'],
    });
  });
});
