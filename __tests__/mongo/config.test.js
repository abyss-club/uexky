import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';
import ConfigModel from '~/models/config';
import { ParamsError } from '~/error';

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

describe('Testing mainTags', () => {
  const mockTags = ['mainA', 'mainB', 'mainC'];
  const modifyTags = ['mainC', 'mainD', 'mainE'];
  const failTags = 'main';
  it('add mainTags', async () => {
    await ConfigModel.modifyMainTags(mockTags);
    const result = await ConfigModel.findOne({ optionName: 'mainTags' });
    expect(result.optionName).toEqual('mainTags');
    expect(JSON.parse(result.optionValue)).toEqual(mockTags);
  });
  it('verify mainTags', async () => {
    const result = await ConfigModel.getMainTags();
    expect(result).toEqual(mockTags);
  });
  it('modify mainTags', async () => {
    await ConfigModel.modifyMainTags(modifyTags);
    const result = await ConfigModel.findOne({ optionName: 'mainTags' });
    expect(result.optionName).toEqual('mainTags');
    expect(JSON.parse(result.optionValue)).toEqual(modifyTags);
  });
  it('verify modified mainTags', async () => {
    const result = await ConfigModel.getMainTags();
    expect(result).toEqual(modifyTags);
  });
  it('add invalid tag string', async () => {
    await expect(ConfigModel.modifyMainTags(failTags)).rejects.toThrow(ParamsError);
  });
});
