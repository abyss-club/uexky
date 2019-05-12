import { ObjectID } from 'mongodb';
import startRepl from '../__utils__/mongoServer';
import mongo from '~/utils/mongo';
import Joi from '@hapi/joi';

import findSlice from '~/models/base';

jest.setTimeout(60000);

let replSet;
let mongoClient;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

const BASE = 'base';

const testSchema = Joi.array().items({
  text: Joi.string(),
  count: Joi.number().integer().min(0),
}).single();

const option = {
  query: { text: 'test string' },
  desc: true,
  field: '_id',
  sliceName: 'test',
  parse: value => new ObjectID(value),
  toCursor: value => value.valueOf(),
};

const altOption = {
  query: { count: { $gt: 2, $lt: 6 } },
  desc: false,
  field: '_id',
  sliceName: 'test',
  parse: value => new ObjectID(value),
  toCursor: value => value.valueOf(),
};

const sliceQuery = { after: '', limit: 10 };

describe('Testing sliceQuery', () => {
  it('Preparation', async () => {
    const { error, value } = Joi.validate([
      { text: 'test string', count: 1 },
      { count: 2 },
      { count: 3 },
      { count: 4 },
      { count: 5 },
      { count: 6 },
    ], testSchema);
    expect(error).toBeNull();
    const r = await mongo.collection(BASE).insertMany(value);
    expect(r.insertedCount).toEqual(6);
  });
  it('Find string', async () => {
    const result = await findSlice(sliceQuery, mongo.collection(BASE), option);
    const count1 = await mongo.collection(BASE).findOne({ count: 1 });
    expect(result.sliceInfo.firstCursor).toEqual(count1._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count1._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(false);
  });
  it('Find between', async () => {
    const result = await findSlice(sliceQuery, mongo.collection(BASE), altOption);
    const count3 = await mongo.collection(BASE).findOne({ count: 3 });
    const count5 = await mongo.collection(BASE).findOne({ count: 5 });
    expect(result.sliceInfo.firstCursor).toEqual(count3._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count5._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(false);
  });
  it('Find between with limitation', async () => {
    const result = await findSlice({ ...sliceQuery, limit: 1 }, mongo.collection(BASE), altOption);
    const count3 = await mongo.collection(BASE).findOne({ count: 3 });
    const count4 = await mongo.collection(BASE).findOne({ count: 4 });
    expect(result.sliceInfo.firstCursor).toEqual(count3._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count4._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(true);
  });
  it('Find between with limitation and desc', async () => {
    const newOption = { ...altOption, desc: true };
    const result = await findSlice(
      { ...sliceQuery, limit: 1 }, mongo.collection(BASE), newOption,
    );
    const count4 = await mongo.collection(BASE).findOne({ count: 4 });
    const count5 = await mongo.collection(BASE).findOne({ count: 5 });
    expect(result.sliceInfo.firstCursor).toEqual(count4._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count5._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(true);
  });
  it('Find using before and after', async () => {
    const count3 = await mongo.collection(BASE).findOne({ count: 3 });
    const count4 = await mongo.collection(BASE).findOne({ count: 4 });
    const count5 = await mongo.collection(BASE).findOne({ count: 5 });
    const result = await findSlice({
      ...sliceQuery, after: count3._id.valueOf(),
    }, mongo.collection(BASE), altOption);
    expect(result.sliceInfo.firstCursor).toEqual(count4._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count5._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(false);
  });
});
