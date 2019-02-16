import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';
import UserModel from '~/models/user';
import { ParamsError, InternalError } from '~/utilities/error';

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

const mockUser = {
  email: 'test@example.com',
  name: 'testUser',
};
const subbedTags = ['MainA', 'SubA', 'SubB', 'MainB'];
const newSubbedTags = [...subbedTags, 'MainC'];
const dupSubbedTags = [...subbedTags, 'MainC', 'MainD'];

it('add user before tests', async () => {
  await UserModel.getUserByEmail(mockUser.email);
});

describe('Testing modifying tags subscription', () => {
  it('sync tags', async () => {
    let user = await UserModel.getUserByEmail(mockUser.email);
    await user.syncTags(subbedTags);
    user = await UserModel.getUserByEmail(mockUser.email);
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(subbedTags));
  });
});
describe('Testing adding tags subscription', () => {
  it('add tags', async () => {
    let user = await UserModel.getUserByEmail(mockUser.email);
    await user.addSubbedTags(['MainC']);
    user = await UserModel.getUserByEmail(mockUser.email);
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(newSubbedTags));
  });
  it('add invalid tags', async () => {
    const user = await UserModel.getUserByEmail(mockUser.email);
    expect(user.addSubbedTags('MainC')).rejects.toThrow(ParamsError);
  });
  it('add duplicated tags', async () => {
    let user = await UserModel.getUserByEmail(mockUser.email);
    await user.addSubbedTags(['MainC', 'MainD']);
    user = await UserModel.getUserByEmail(mockUser.email);
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(dupSubbedTags));
  });
});
describe('Testing deleting tags subscription', () => {
  it('del tags', async () => {
    let user = await UserModel.getUserByEmail(mockUser.email);
    await user.delSubbedTags(['MainD']);
    user = await UserModel.getUserByEmail(mockUser.email);
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(newSubbedTags));
  });
  it('del invalid tags', async () => {
    const user = await UserModel.getUserByEmail(mockUser.email);
    expect(user.delSubbedTags('MainC')).rejects.toThrow(ParamsError);
  });
  it('del non-existing tags', async () => {
    let user = await UserModel.getUserByEmail(mockUser.email);
    await user.delSubbedTags(['MainC', 'MainD']);
    user = await UserModel.getUserByEmail(mockUser.email);
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(subbedTags));
  });
});
describe('Testing setting name', () => {
  it('set name', async () => {
    let user = await UserModel.getUserByEmail(mockUser.email);
    await user.setName(mockUser.name);
    user = await UserModel.getUserByEmail(mockUser.email);
    expect(user.name).toEqual(mockUser.name);
  });
  it('set name again', async () => {
    const user = await UserModel.getUserByEmail(mockUser.email);
    expect(user.setName(mockUser.name)).rejects.toThrow(InternalError);
  });
});
