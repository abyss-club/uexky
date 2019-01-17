import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';
import ConfigModel from '~/models/config';
import { InternalError, ParamsError } from '~/error';

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

describe('Testing rateLimit', () => {
  const emptySettings = {};
  const validSettings = {
    HTTPHeader: 'Header',
    QueryLimit: 1,
    QueryResetTime: 2,
    MutLimit: 3,
    MutResetTime: 4,
    Cost: {
      CreateUser: 5,
      PubThread: 6,
      PubPost: 7,
    },
  };
  const incompleteSettings = {
    HTTPHeader: 'Modified',
    QueryLimit: 10,
    QueryResetTime: 20,
    MutLimit: 30,
    MutResetTime: 40,
  };
  const invalidSettings = {
    HTTPHeader: 0,
    QueryLimit: '1',
    QueryResetTime: '2',
    MutLimit: '3',
    MutResetTime: '4',
    Cost: {
      CreateUser: '5',
      PubThread: '6',
      PubPost: '7',
    },
  };
  const finalSettings = {
    HTTPHeader: 'Modified',
    QueryLimit: 10,
    QueryResetTime: 20,
    MutLimit: 30,
    MutResetTime: 40,
    Cost: {
      CreateUser: 30,
      PubThread: 10,
      PubPost: 1,
    },
  };
  const defaultSettings = {
    HTTPHeader: '',
    QueryLimit: 300,
    QueryResetTime: 3600,
    MutLimit: 30,
    MutResetTime: 3600,
    Cost: {
      CreateUser: 30,
      PubThread: 10,
      PubPost: 1,
    },
  };
  it('add valid settings', async () => {
    const returned = await ConfigModel.modifyRateLimit(validSettings);
    const result = await ConfigModel.findOne({ optionName: 'rateLimit' });
    expect(returned).toEqual(JSON.stringify(validSettings));
    expect(JSON.parse(result.optionValue)).toEqual(validSettings);
  });
  it('verify valid settings', async () => {
    const result = await ConfigModel.getRateLimit();
    expect(result).toEqual(JSON.stringify(validSettings));
  });
  it('add empty settings', async () => {
    await mongoose.connection.db.dropDatabase();
    const returned = await ConfigModel.modifyRateLimit(emptySettings);
    const result = await ConfigModel.findOne({ optionName: 'rateLimit' });
    expect(returned).toEqual(JSON.stringify(defaultSettings));
    expect(JSON.parse(result.optionValue)).toEqual(defaultSettings);
  });
  it('verify empty settings', async () => {
    const result = await ConfigModel.getRateLimit();
    expect(result).toEqual(JSON.stringify(defaultSettings));
  });
  it('add invalid settings', async () => {
    await mongoose.connection.db.dropDatabase();
    await expect(ConfigModel.modifyRateLimit(invalidSettings)).rejects.toThrow(ParamsError);
  });
  it('verify invalid settings', async () => {
    expect(ConfigModel.getRateLimit()).rejects.toThrow(InternalError);
  });
  it('add incomplete settings', async () => {
    const returned = await ConfigModel.modifyRateLimit(incompleteSettings);
    const result = await ConfigModel.findOne({ optionName: 'rateLimit' });
    expect(returned).toEqual(JSON.stringify(finalSettings));
    expect(JSON.parse(result.optionValue)).toEqual(finalSettings);
  });
  it('verify incomplete settings', async () => {
    const result = await ConfigModel.getRateLimit();
    expect(result).toEqual(JSON.stringify(finalSettings));
  });
  it('add valid, then incomplete settings', async () => {
    await mongoose.connection.db.dropDatabase();
    const validReturn = await ConfigModel.modifyRateLimit(validSettings);
    const returned = await ConfigModel.modifyRateLimit(incompleteSettings);
    const result = await ConfigModel.findOne({ optionName: 'rateLimit' });
    expect(validReturn).toEqual(JSON.stringify(validSettings));
    expect(returned).toEqual(JSON.stringify(finalSettings));
    expect(JSON.parse(result.optionValue)).toEqual(finalSettings);
  });
  it('verify final settings', async () => {
    const result = await ConfigModel.getRateLimit();
    expect(result).toEqual(JSON.stringify(finalSettings));
  });
});
