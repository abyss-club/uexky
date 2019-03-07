import gql from 'graphql-tag';
import mongoose from 'mongoose';
import { startMongo } from '../__utils__/mongoServer';

import { mutate } from '../__utils__/apolloClient';
import TagModel from '~/models/tag';

let mongoServer;

beforeAll(async () => {
  mongoServer = await startMongo();
});

afterAll(() => {
  mongoose.disconnect();
  mongoServer.stop();
});

const mockTags = { mainTag: 'MainA', subTags: ['SubA', 'SubB'] };

// const ADD_TAGS = gql`
//   mutation AddTags($tags: [String!]) {
//     editConfig(config: {
//       mainTags: $tags,
//     }) { mainTags }
//   }
// `;

const GET_TAGS = gql`
  query {
    tags {
      mainTags,
      tree {
        mainTag, subTags
      }
    }
  }
`;

describe('Insert Tags', () => {
  it('add tags', async () => {
    await TagModel.addMainTag(mockTags.mainTag);
    const tags = mockTags;
    await TagModel.onPubThread(tags);
  });
  it('validate tags', async () => {
    const result = await mutate({ query: GET_TAGS });
    expect(result.data.tags.mainTags).toEqual([mockTags.mainTag]);
    const target = result.data.tags.tree.filter(tagObj => tagObj.mainTag === mockTags.mainTag)[0];
    expect(target.mainTag).toEqual(mockTags.mainTag);
    expect(target.subTags.sort()).toEqual(mockTags.subTags.sort());
  });
});
