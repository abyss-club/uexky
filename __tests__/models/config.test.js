import mongoose from 'mongoose';
import ConfigModel from '~/models/config';
import { ParamsError } from '~/utils/error';

import { startMongo } from '../__utils__/mongoServer';

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
  const expectConfig = {
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
    const result = await ConfigModel.getConfig();
    const resultInDb = await ConfigModel.findOne().exec();
    expect(result).toEqual(expectConfig);
    expect(resultInDb.format()).toEqual(expectConfig);
  };

  it('get default config', async () => {
    const result = await ConfigModel.getConfig();
    expect(result).toEqual(expectConfig);
  });
  it('modify config with single entry', async () => {
    await ConfigModel.setConfig({ rateLimit: { httpHeader: 'X-IP-Forward' } });
    expectConfig.rateLimit.httpHeader = 'X-IP-Forward';
    await checkConfig(expectConfig);
  });
  it('modify config with invalid value (unknown group)', async () => {
    await expect(ConfigModel.setConfig({
      iamfine: true,
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig(expectConfig);
  });
  it('modify config with invalid value (unknown entry)', async () => {
    await expect(ConfigModel.setConfig({
      rateLimit: { mutLimit: 'hello', name: 'tom' },
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig(expectConfig);
  });
  it('modify config with invalid value (error group)', async () => {
    await expect(ConfigModel.setConfig({
      rateLimit: 'unreal',
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig(expectConfig);
  });
  it('modify config with invalid value (error entry)', async () => {
    await expect(ConfigModel.setConfig({
      rateLimit: { queryLimit: 'hi' },
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig(expectConfig);
  });
  it('modify config multiple entries', async () => {
    await ConfigModel.setConfig({
      rateLimit: { queryLimit: 400 },
      rateCost: { pubThread: 20, pubPost: 3 },
    });
    expectConfig.rateLimit.queryLimit = 400;
    expectConfig.rateCost.pubThread = 20;
    expectConfig.rateCost.pubPost = 3;
    await checkConfig(expectConfig);
  });
});
