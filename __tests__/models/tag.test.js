import mongoose from 'mongoose';
import { startMongo } from '../__utils__/mongoServer';
import TagModel from '~/models/tag';

// May require additional time for downloading MongoDB binaries
// jasmine.DEFAULT_TIMEOUT_INTERVAL = 600000;

let mongoServer;

beforeAll(async () => {
  mongoServer = await startMongo();
});

afterAll(() => {
  mongoose.disconnect();
  mongoServer.stop();
});
const mockTags = { mainTag: 'MainA', subTags: ['SubA', 'SubB'] };

describe('Insert Tags', () => {
  it('add tags', async () => {
    await TagModel.addMainTag(mockTags.mainTag);
    await TagModel.onPubThread(mockTags);
  });
  it('validate tags', async () => {
    const result = await TagModel.getTree();
    const target = result.filter(tagObj => tagObj.mainTag === mockTags.mainTag)[0];
    expect(target.mainTag).toEqual(mockTags.mainTag);
    expect(target.subTags.sort()).toEqual(mockTags.subTags.sort());
  });
});
