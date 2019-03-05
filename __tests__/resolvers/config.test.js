import gql from 'graphql-tag';
import mongoose from 'mongoose';
import { startMongo } from '../__utils__/mongoServer';

import { query, mutate } from '../__utils__/apolloClient';

let mongoServer;

beforeAll(async () => {
  mongoServer = await startMongo();
});

afterAll(() => {
  mongoose.disconnect();
  mongoServer.stop();
});

// const EDIT_TAGS = gql`
//   mutation AddTags($tags: [String!]) {
//     editConfig(config: {
//       mainTags: $tags,
//     }) { mainTags }
//   }
// `;

const EDIT_CONFIG = gql`
  mutation editConfig($config: ConfigInput!) {
    editConfig(config: $config) {
      rateLimit {
        httpHeader, queryLimit, queryResetTime, mutLimit, mutResetTime
      }
      rateCost {
        createUser, pubThread, pubPost
      }
    }
  }
`;

const GET_RATE = gql`
  query {
    config {
      rateLimit {
        httpHeader, queryLimit, queryResetTime, mutLimit, mutResetTime
      }
      rateCost {
        createUser, pubThread, pubPost
      }
    }
  }
`;

// describe('Testing mainTags', () => {
//   const mockTags = ['mainA', 'mainB', 'mainC'];
//   const modifyTags = ['mainC', 'mainD', 'mainE'];
//   const failTags = '';
//   it('add mainTags', async () => {
//     const result = await mutate({ mutation: EDIT_TAGS, variables: { tags: mockTags } });
//     expect(result.data.editConfig.mainTags).toEqual(mockTags);
//   });
//   it('modify mainTags', async () => {
//     const result = await mutate({ mutation: EDIT_TAGS, variables: { tags: modifyTags } });
//     expect(result.data.editConfig.mainTags).toEqual(modifyTags);
//   });
//   it('add invalid tag string', async () => {
//     const result = await mutate({ mutation: EDIT_TAGS, variables: { tags: failTags } });
//     expect(result.data).toBeNull();
//     expect(result.errors[0].message).toEqual('Invalid tag provided in array.');
//   });
// });

describe('Testing rateLimit', () => {
  const expectedConfig = {
    rateLimit: {
      httpHeader: '',
      queryLimit: 300,
      queryResetTime: 3600,
      mutLimit: 30,
      mutResetTime: 3600,
    },
    rateCost: {
      createUser: 30,
      pubThread: 10,
      pubPost: 1,
    },
  };

  it('get default config', async () => {
    const result = await query({ query: GET_RATE });
    expect(result.data.config).toEqual(expectedConfig);
  });
  it('modify config with single entry', async () => {
    const result = await mutate({
      mutation: EDIT_CONFIG,
      variables: {
        config: {
          rateLimit: { httpHeader: 'X-IP-Forward' },
        },
      },
    });
    expectedConfig.rateLimit.httpHeader = 'X-IP-Forward';
    expect(result.data.editConfig).toEqual(expectedConfig);
  });
  it('modify config with invalid httpHeader', async () => {
    const result = await mutate({
      mutation: EDIT_CONFIG,
      variables: {
        config: {
          rateLimit: { httpHeader: ';?:""(&[])' },
        },
      },
    });
    expect(result.errors.length).toEqual(1);
  });
  it('modify config multiple entries', async () => {
    const result = await mutate({
      mutation: EDIT_CONFIG,
      variables: {
        config: {
          rateLimit: { queryLimit: 400 },
          rateCost: { pubThread: 20, pubPost: 3 },
        },
      },
    });
    expectedConfig.rateLimit.queryLimit = 400;
    expectedConfig.rateCost.pubThread = 20;
    expectedConfig.rateCost.pubPost = 3;
    expect(result.data.editConfig).toEqual(expectedConfig);
  });
});
