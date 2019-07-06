import UID, { Base64 } from '~/uid';
import startPg, { migrate } from '../__utils__/pgServer';

describe('Base64', () => {
  const pairs = [
    // 2000 = 31 * 64 + 16
    [BigInt(2000), Base64.code[31] + Base64.code[16]],
    // 200 = 3 * 64 + 8
    [BigInt(200), Base64.code[3] + Base64.code[8]],
    // 2 = 0 * 64 + 2
    [BigInt(2), Base64.code[0] + Base64.code[2]],
  ];
  test('convertToBigInt', () => {
    pairs.forEach((p) => {
      const n = Base64.convertToBigInt(p[1]);
      expect(n).toEqual(p[0]);
    });
  });
  test('parseFromBigInt', () => {
    pairs.forEach((p) => {
      const s = Base64.parseFromBigInt(p[0], 2);
      expect(s).toEqual(p[1]);
    });
  });
});

describe('UID', () => {
  let pgPool;
  beforeAll(async () => {
    await migrate();
    pgPool = await startPg();
  });
  afterAll(async () => {
    await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
    pgPool.end();
  });
  test('new id', async () => {
    const uid = await UID.new();
    expect(typeof uid.suid).toEqual('bigint');
    expect(uid.duid).toMatch(/[0-9a-zA-Z-_]{6,}/);
  });
  test('parse duid', async () => {
    const uid = await UID.new();
    const uid2 = UID.parse(uid.duid);
    expect(uid2.duid).toEqual(uid.duid);
    expect(uid2.suid.toString()).toEqual(uid.suid.toString());
  });
  test('parse suid', async () => {
    const uid = await UID.new();
    const uid3 = UID.parse(uid.suid);
    expect(uid3.duid).toEqual(uid.duid);
    expect(uid3.suid.toString()).toEqual(uid.suid.toString());
  });
  test('parse uid from uid', async () => {
    const uid = await UID.new();
    const uid3 = UID.parse(uid);
    expect(uid3.duid).toEqual(uid.duid);
    expect(uid3.suid.toString()).toEqual(uid.suid.toString());
  });
});
