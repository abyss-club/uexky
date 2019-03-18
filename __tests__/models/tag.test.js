import { startRepl } from '../__utils__/mongoServer';
import TagModel from '~/models/tag';

jest.setTimeout(60000);

let replSet;
let mongoClient;
// let db;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

// const TAG = 'tag';

const mockTags = { mainTag: 'MainA', subTags: ['SubA', 'SubB'] };

describe('Insert Tags', () => {
  it('add tags', async () => {
    await TagModel().addMainTag(mockTags.mainTag);
    await TagModel().onPubThread(mockTags);
  });
  it('validate tags', async () => {
    const result = await TagModel().getTree();
    const target = result.filter(tagObj => tagObj.mainTag === mockTags.mainTag)[0];
    expect(target.mainTag).toEqual(mockTags.mainTag);
    expect(target.subTags.sort()).toEqual(mockTags.subTags.sort());
  });
});
