import ConfigModel from '~/models/config';
import { ROLE } from '~/models/user';
import { ParamsError } from '~/utils/error';
import startPg, { migrate } from '../__utils__/pgServer';
import mockContext from '../__utils__/context';

let pgPool;
// let db;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

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
    const result = await ConfigModel.getConfig();
    const resultInDb = await pgPool.query(
      'SELECT "rateLimit", "rateCost" from config where id = 1',
    );
    expect(result).toEqual(expectedConfig);
    expect(resultInDb.rows[0]).toEqual(expectedConfig);
  };
  let ctx;
  it('prepare', async () => {
    ctx = await mockContext({ email: 'uexky@uexky.com', role: ROLE.ADMIN });
  });
  it('get default config', async () => {
    const result = await ConfigModel.getConfig();
    expect(result).toEqual(expectedConfig);
  });
  it('modify config with single entry', async () => {
    await ConfigModel.setConfig(ctx, { rateLimit: { httpHeader: 'X-IP-Forward' } });
    expectedConfig.rateLimit.httpHeader = 'X-IP-Forward';
    await checkConfig();
  });
  it('modify config with invalid value (unknown group)', async () => {
    await expect(ConfigModel.setConfig(ctx, {
      iamfine: true,
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig();
  });
  it('modify config with invalid value (unknown entry)', async () => {
    await expect(ConfigModel.setConfig(ctx, {
      rateLimit: { mutLimit: 'hello', name: 'tom' },
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig();
  });
  it('modify config with invalid value (error group)', async () => {
    await expect(ConfigModel.setConfig(ctx, {
      rateLimit: 'unreal',
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig();
  });
  it('modify config with invalid value (error entry)', async () => {
    await expect(ConfigModel.setConfig(ctx, {
      rateLimit: { queryLimit: 'hi' },
      rateCost: { pubPost: 2 },
    })).rejects.toThrow(ParamsError);
    await checkConfig();
  });
  it('modify config multiple entries', async () => {
    await ConfigModel.setConfig(ctx, {
      rateLimit: { queryLimit: 400 },
      rateCost: { pubThread: 20, pubPost: 3 },
    });
    expectedConfig.rateLimit.queryLimit = 400;
    expectedConfig.rateCost.pubThread = 20;
    expectedConfig.rateCost.pubPost = 3;
    await checkConfig();
  });
});
