import TagModel from '~/models/tag';
import ThreadModel from '~/models/thread';
import { query } from '~/utils/pg';
import startPg, { migrate } from '../__utils__/pgServer';
import mockContext from '../__utils__/context';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

describe('tags query', () => {
  it('parpare data', async () => {
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['MainA', true]);
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['MainB', true]);
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['SubA', false]);
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['SubB', false]);
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['SubC', false]);
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['SubD', false]);
  });
  it('find one tag', async () => {
    const tags = await TagModel.findTags({ query: 'SubA' });
    expect(tags.length).toEqual(1);
    expect(tags[0].name).toEqual('SubA');
    expect(tags[0].isMain).toBeFalsy();
  });
  it('find tags by query string', async () => {
    const tags = await TagModel.findTags({ query: 'B' });
    expect(tags.length).toEqual(2);
    expect(tags[0].name).toEqual('SubB');
    expect(tags[0].isMain).toBeFalsy();
    expect(tags[1].name).toEqual('MainB');
    expect(tags[1].isMain).toBeTruthy();
  });
  it('find main tags', async () => {
    const mainTags = await TagModel.getMainTags();
    expect(mainTags).toEqual(['MainA', 'MainB']);
  });
});

describe('set thread tags', () => {
  let thread;
  it('parpare data', async () => {
    await query('DELETE FROM tag');
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['MainA', true]);
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['MainB', true]);
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['SubA', false]);
    const ctx = await mockContext({ email: 'test@uexky.com' });
    thread = await ThreadModel.new({
      ctx,
      thread: {
        anonymous: true,
        content: 'Test Content',
        mainTag: 'MainA',
        subTags: ['SubA', 'SubB'],
        title: 'TestTitle',
      },
    });
  });
  it('check in tag', async () => {
    const tags = await TagModel.findTags({ }); // all
    expect(tags.length).toEqual(4);
    expect(tags[3].name).toEqual('MainB'); // oldest
  });
  it('check threads_tags', async () => {
    const { rows } = await query(
      'SELECT "tagName" as name FROM threads_tags WHERE "threadId"=$1',
      [thread.id.suid],
    );
    const tags = rows.map(row => row.name);
    expect(tags.length).toEqual(3);
    expect(tags).toContain('MainA');
    expect(tags).toContain('SubA');
    expect(tags).toContain('SubB');
  });
  it('check tags_main_tags', async () => {
    const { rows } = await query(
      'SELECT name, "belongsTo" FROM tags_main_tags WHERE "belongsTo"=$1',
      ['MainA'],
    );
    expect(rows.length).toEqual(2);
    expect(rows).toContainEqual({ name: 'SubA', belongsTo: 'MainA' });
    expect(rows).toContainEqual({ name: 'SubB', belongsTo: 'MainA' });
  });
});
