import UserModel from '~/models/user';
import { query } from '~/utils/pg';
import { ParamsError, InternalError } from '~/utils/error';
import startPg, { migrate } from '../__utils__/pgServer';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

describe('user object', () => {
  const mockEmail = 'test1@example.com';
  it('new user context', async () => {
    const auth = await UserModel.authContext({ email: mockEmail });
    const user = auth.signedInUser();
    expect(user.email).toEqual(mockEmail);
  });
  it('user context', async () => {
    const auth = await UserModel.authContext({ email: mockEmail });
    const user = auth.signedInUser();
    expect(user.email).toEqual(mockEmail);
    const { rows } = await query('SELECT count(*) FROM public.user');
    expect(parseInt(rows[0].count, 10)).toEqual(1);
  });
  it('find user', async () => {
    const user = await UserModel.findByEmail({ email: mockEmail });
    expect(user.email).toEqual(mockEmail);
  });
});

describe('user name', () => {
  const mockEmail = 'test1@example.com';
  const mockName = 'test1';
  it('set and get name', async () => {
    const auth = await UserModel.authContext({ email: mockEmail });
    const ctx = { auth };
    const user = await UserModel.setName({ ctx, name: mockName });
    expect(user.name).toEqual(mockName);
    const userInDb = await UserModel.findByEmail({ email: mockEmail });
    expect(userInDb.name).toEqual(mockName);
  });
  it('name is unchangeable', async () => {
    const auth = await UserModel.authContext({ email: mockEmail });
    const ctx = { auth };
    await expect(UserModel.setName({ ctx, name: mockName })).rejects.toThrow(ParamsError);
  });
  it('name is unique', async () => {
    const auth = await UserModel.authContext({ email: 'test2@example.com' });
    const ctx = { auth };
    await expect(UserModel.setName({ ctx, name: mockName })).rejects.toThrow(InternalError);
  });
});

describe('user tags', () => {
  const tagsInDb = [
    { name: 'MainA', isMain: true },
    { name: 'MainB', isMain: true },
    { name: 'MainC', isMain: true },
    { name: 'SubA', isMain: false },
    { name: 'SubB', isMain: false },
    { name: 'SubC', isMain: false },
    { name: 'SubD', isMain: false },
  ];
  const mockEmail = 'test@example.com';
  let ctx;
  it('build data', async () => {
    await Promise.all(tagsInDb.map(tag => query(
      'INSERT INTO tag (name, "isMain") VALUES ($1, $2)',
      [tag.name, tag.isMain],
    )));
    const auth = await UserModel.authContext({ email: mockEmail });
    ctx = { auth };
  });
  it('add tag', async () => {
    await UserModel.addSubbedTag({ ctx, tag: 'MainA' });
    await UserModel.addSubbedTag({ ctx, tag: 'SubA' });
    const user = await UserModel.findByEmail({ email: mockEmail });
    const tags = await user.getTags();
    expect(tags.length).toEqual(2);
    expect(tags).toContain('MainA');
    expect(tags).toContain('SubA');
  });
  it('add tag invalid', async () => {
    await expect(UserModel.addSubbedTag({ ctx, tag: 'SubA' })).rejects.toThrow(InternalError);
    await expect(UserModel.addSubbedTag({ ctx, tag: 'SubX' })).rejects.toThrow(InternalError);
  });
  it('del tag', async () => {
    await UserModel.delSubbedTag({ ctx, tag: 'MainA' });
    const user = await UserModel.findByEmail({ email: mockEmail });
    const tags = await user.getTags();
    expect(tags.length).toEqual(1);
    expect(tags).toContain('SubA');
  });
  it('del tag invalid', async () => {
    await UserModel.delSubbedTag({ ctx, tag: 'SubB' });
    await UserModel.delSubbedTag({ ctx, tag: 'SubX' });
    const user = await UserModel.findByEmail({ email: mockEmail });
    const tags = await user.getTags();
    expect(tags.length).toEqual(1);
    expect(tags).toContain('SubA');
  });
  it('sync tag', async () => {
    await UserModel.syncTags({ ctx, tags: ['MainC', 'SubC', 'SubD'] });
    const user = await UserModel.findByEmail({ email: mockEmail });
    const tags = await user.getTags();
    expect(tags.length).toEqual(3);
    expect(tags).toContain('MainC');
    expect(tags).toContain('SubC');
    expect(tags).toContain('SubD');
  });
  it('sync tag invalid', async () => {
    await expect(UserModel.syncTags({ ctx, tags: ['SubC', 'SubX'] })).rejects.toThrow(InternalError);
    const user = await UserModel.findByEmail({ email: mockEmail });
    const tags = await user.getTags();
    expect(tags.length).toEqual(3);
    expect(tags).toContain('MainC');
    expect(tags).toContain('SubC');
    expect(tags).toContain('SubD');
  });
});
