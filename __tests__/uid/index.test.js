import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';

import Uid, { Base64 } from '~/uid';

let mongoServer;
const opts = { useNewUrlParser: true };
const connectMongodb = async () => {
  mongoServer = new MongoMemoryServer();
  const mongoUri = await mongoServer.getConnectionString();
  await mongoose.connect(mongoUri, opts, (err) => {
    if (err) console.error(err);
  });
};
const disconnectMongodb = () => {
  mongoose.disconnect();
  mongoServer.stop();
};

describe('Base64', () => {
  const pairs = [
    // 2000 = 31 * 64 + 16 = 0x7d0
    ['7d0', Base64.code[31] + Base64.code[16]],
    // 200 = 3 * 64 + 8 = 0xc8
    ['0c8', Base64.code[3] + Base64.code[8]],
    // 2 = 0 * 64 + 2 = 0x2
    ['002', Base64.code[0] + Base64.code[2]],
  ];
  test('convertTohex3', () => {
    pairs.forEach((p) => {
      const hex = Base64.convertTohex3(p[1]);
      expect(hex).toEqual(p[0]);
    });
  });
  test('parseFromHex3', () => {
    pairs.forEach((p) => {
      const b = Base64.parseFromHex3(p[0]);
      expect(b).toEqual(p[1]);
    });
  });
});

describe('decode/encode', () => {
  beforeEach(connectMongodb);
  afterEach(disconnectMongodb);
  test('Generator id, decode, and encode', async () => {
    const suid = await Uid.newSuid();
    expect(suid).toMatch(/[0-9a-f]{15}/);
    const uid = Uid.decode(suid);
    expect(uid).toMatch(/[0-9a-zA-Z-_]{10}/);
    const suid1 = Uid.encode(uid);
    expect(suid1).toEqual(suid);
  });
});

test('timestamp reverse', () => {
  const suid = '800000000000000';
  // 1000-0000-0000-0000-0000-0000-0000-0000 = 80000000
  // 000000-000000-000000-000000-0000000-010000 = AAAAAg
  const uid = Uid.decode(suid);
  expect(uid.substring(0, 6)).toEqual('AAAAAQ');
});
