import ConfigModel from '~/models/config';
import dbClient from '~/dbClient';
import { ParamsError } from '~/utils/error';
import { startRepl } from '../__utils__/mongoServer';

// May require additional time for downloading MongoDB binaries
jest.setTimeout(60000);

let replSet;
let mongoClient;
let ctx;
// let db;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

const CONFIG = 'config';

/* TODO(tangwenhan): may be useful in next work (tag model).
describe('Testing mainTags', () => {
  const mockTags = ['mainA', 'mainB', 'mainC'];
  const modifyTags = ['mainC', 'mainD', 'mainE'];
  const failTags = 'main';
  it('add mainTags', async () => {
    await ConfigModel.modifyMainTags(mockTags);
    const result = await ConfigModel.findOne({ optionName: 'mainTags' }).exec();
    expect(result.optionName).toEqual('mainTags');
    expect(JSON.parse(result.optionValue)).toEqual(mockTags);
  });
  it('verify mainTags', async () => {
    const result = await ConfigModel.getMainTags();
    expect(result).toEqual(mockTags);
  });
  it('modify mainTags', async () => {
    await ConfigModel.modifyMainTags(modifyTags);
    const result = await ConfigModel.findOne({ optionName: 'mainTags' }).exec();
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
*/

describe('Testing rateLimit', () => {
  // default
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
  const checkConfig = async () => {
    const result = await ConfigModel().getConfig();
    const resultInDb = await dbClient.collection(CONFIG).findOne({}, { projection: { _id: 0 } });
    expect(result).toEqual(expectedConfig);
    expect(resultInDb).toEqual(expectedConfig);
  };

  it('get default config', async () => {
    const result = await ConfigModel().getConfig();
    expect(result).toEqual(expectedConfig);
  });
  it('modify config with single entry', async () => {
    await ConfigModel(ctx).setConfig({ rateLimit: { httpHeader: 'X-IP-Forward' } });
    expectedConfig.rateLimit.httpHeader = 'X-IP-Forward';
    await checkConfig();
  });
  it('modify config with invalid value (unknown group)', async () => {
    await expect(ConfigModel(ctx).setConfig({
      iamfine: true,
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig();
  });
  it('modify config with invalid value (unknown entry)', async () => {
    await expect(ConfigModel(ctx).setConfig({
      rateLimit: { mutLimit: 'hello', name: 'tom' },
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig();
  });
  it('modify config with invalid value (error group)', async () => {
    await expect(ConfigModel(ctx).setConfig({
      rateLimit: 'unreal',
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig();
  });
  it('modify config with invalid value (error entry)', async () => {
    await expect(ConfigModel(ctx).setConfig({
      rateLimit: { queryLimit: 'hi' },
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig();
  });
  it('modify config multiple entries', async () => {
    await ConfigModel(ctx).setConfig({
      rateLimit: { queryLimit: 400 },
      rateCost: { pubThread: 20, pubPost: 3 },
    });
    expectedConfig.rateLimit.queryLimit = 400;
    expectedConfig.rateCost.pubThread = 20;
    expectedConfig.rateCost.pubPost = 3;
    await checkConfig();
  });
});
