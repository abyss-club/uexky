import fs from 'fs';
import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';
import TagModel from '~/models/tag';
import ThreadModel from '~/models/thread';

// May require additional time for downloading MongoDB binaries
// jasmine.DEFAULT_TIMEOUT_INTERVAL = 600000;

let mongoServer;
const opts = { useNewUrlParser: true };

beforeAll(async () => {
  mongoServer = new MongoMemoryServer();
  const mongoUri = await mongoServer.getConnectionString();
  await mongoose.connect(mongoUri, opts, (err) => {
    if (err) console.error(err);
  });
});

afterAll(() => {
  mongoose.disconnect();
  mongoServer.stop();
});

describe('Insert Tags', () => {
  const mockTags = { mainTag: 'MainA', subTags: ['SubA', 'SubB'] };
  it('add tags', async () => {
    const tags = mockTags;
    // session.startTransaction();

    await TagModel.onPubThread(tags);
    // await session.commitTransaction();
  });
  it('validate tags', async () => {
    const result = await TagModel.getTree();
    const target = result.filter(tagObj => tagObj.mainTag === mockTags.mainTag)[0];
    expect(target.mainTag).toEqual(mockTags.mainTag);
    expect(target.subTags).toEqual(mockTags.subTags);
  });
});
