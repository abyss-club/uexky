import startRepl from '../__utils__/mongoServer';

import UserModel from '~/models/user';
import { ParamsError, InternalError } from '~/utils/error';

jest.setTimeout(60000); // for boot replica sets
let replSet;
let mongoClient;
// let db;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

const mockUser = {
  email: 'test@example.com',
  name: 'testUser',
};
const subbedTags = ['MainA', 'SubA', 'SubB', 'MainB'];
const newSubbedTags = [...subbedTags, 'MainC'];
const dupSubbedTags = [...subbedTags, 'MainC', 'MainD'];

it('add user before tests', async () => {
  await UserModel().getUserByEmail(mockUser.email);
});

describe('Testing modifying tags subscription', () => {
  it('sync tags', async () => {
    let user = await UserModel().getUserByEmail(mockUser.email);
    const result = await UserModel({ user }).methods(user).syncTags(subbedTags);
    user = await UserModel().getUserByEmail(mockUser.email);
    expect(JSON.stringify(result.tags)).toEqual(JSON.stringify(subbedTags));
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(subbedTags));
  });
});
describe('Testing adding tags subscription', () => {
  it('add tags', async () => {
    let user = await UserModel().getUserByEmail(mockUser.email);
    const result = await UserModel({ user }).methods(user).addSubbedTags(['MainC']);
    user = await UserModel().getUserByEmail(mockUser.email);
    expect(JSON.stringify(result.tags)).toEqual(JSON.stringify(newSubbedTags));
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(newSubbedTags));
  });
  it('add invalid tags', async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    expect(UserModel({ user }).methods(user).addSubbedTags('MainC')).rejects.toThrow(ParamsError);
  });
  it('add duplicated tags', async () => {
    let user = await UserModel().getUserByEmail(mockUser.email);
    const result = await UserModel({ user }).methods(user).addSubbedTags(['MainC', 'MainD']);
    user = await UserModel().getUserByEmail(mockUser.email);
    expect(JSON.stringify(result.tags)).toEqual(JSON.stringify(dupSubbedTags));
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(dupSubbedTags));
  });
});
describe('Testing deleting tags subscription', () => {
  it('del tags', async () => {
    let user = await UserModel().getUserByEmail(mockUser.email);
    const result = await UserModel({ user }).methods(user).delSubbedTags(['MainD']);
    user = await UserModel().getUserByEmail(mockUser.email);
    expect(JSON.stringify(result.tags)).toEqual(JSON.stringify(newSubbedTags));
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(newSubbedTags));
  });
  it('del invalid tags', async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    expect(UserModel({ user }).methods(user).delSubbedTags('MainC')).rejects.toThrow(ParamsError);
  });
  it('del non-existing tags', async () => {
    let user = await UserModel().getUserByEmail(mockUser.email);
    const result = await UserModel({ user }).methods(user).delSubbedTags(['MainC', 'MainD']);
    user = await UserModel().getUserByEmail(mockUser.email);
    expect(JSON.stringify(result.tags)).toEqual(JSON.stringify(subbedTags));
    expect(JSON.stringify(user.tags)).toEqual(JSON.stringify(subbedTags));
  });
});
describe('Testing setting name', () => {
  it('set name', async () => {
    let user = await UserModel().getUserByEmail(mockUser.email);
    const result = await UserModel({ user }).methods(user).setName(mockUser.name);
    user = await UserModel().getUserByEmail(mockUser.email);
    expect(result.name).toEqual(mockUser.name);
    expect(user.name).toEqual(mockUser.name);
  });
  it('set name again', async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    expect(UserModel({ user }).methods(user).setName(mockUser.name)).rejects.toThrow(InternalError);
  });
});
