import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';

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

describe('...', () => {
  it('...', async () => {
    const User = mongoose.model('User', new mongoose.Schema({ name: String }));
    const cnt = await User.count();
    expect(cnt).toEqual(0);
  });
});
