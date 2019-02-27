import mongoose from 'mongoose';
import { startMongo } from '../__utils__/mongoServer';

import findSlice from '~/models/base';

// May require additional time for downloading MongoDB binaries
// jasmine.DEFAULT_TIMEOUT_INTERVAL = 600000;

let mongoServer;

beforeAll(async () => {
  mongoServer = await startMongo();
});

afterAll(() => {
  mongoose.disconnect();
  mongoServer.stop();
});

const { ObjectId } = mongoose.Types;

const TestSchema = new mongoose.Schema({
  text: String,
  count: Number,
});

const TestModel = mongoose.model('Test', TestSchema);

const option = {
  query: { text: 'test string' },
  desc: true,
  field: '_id',
  sliceName: 'test',
  parse: value => ObjectId(value),
  toCursor: value => value.valueOf(),
};

const altOption = {
  query: { count: { $gt: 2, $lt: 6 } },
  desc: false,
  field: '_id',
  sliceName: 'test',
  parse: value => ObjectId(value),
  toCursor: value => value.valueOf(),
};

const sliceQuery = { after: '', limit: 10 };

describe('Testing sliceQuery', () => {
  it('Preparation', async () => {
    await TestModel.create({ text: 'test string', count: 1 });
    await TestModel.create({ count: 2 });
    await TestModel.create({ count: 3 });
    await TestModel.create({ count: 4 });
    await TestModel.create({ count: 5 });
    await TestModel.create({ count: 6 });
  });
  it('Find string', async () => {
    const result = await findSlice(sliceQuery, TestModel, option);
    const count1 = await TestModel.findOne({ count: 1 }).exec();
    expect(result.sliceInfo.firstCursor).toEqual(count1._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count1._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(false);
  });
  it('Find between', async () => {
    const result = await findSlice(sliceQuery, TestModel, altOption);
    const count3 = await TestModel.findOne({ count: 3 }).exec();
    const count5 = await TestModel.findOne({ count: 5 }).exec();
    expect(result.sliceInfo.firstCursor).toEqual(count3._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count5._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(false);
  });
  it('Find between with limitation', async () => {
    const result = await findSlice({ ...sliceQuery, limit: 1 }, TestModel, altOption);
    const count3 = await TestModel.findOne({ count: 3 }).exec();
    const count4 = await TestModel.findOne({ count: 4 }).exec();
    expect(result.sliceInfo.firstCursor).toEqual(count3._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count4._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(true);
  });
  it('Find between with limitation and desc', async () => {
    const newOption = { ...altOption, desc: true };
    const result = await findSlice({ ...sliceQuery, limit: 1 }, TestModel, newOption);
    const count4 = await TestModel.findOne({ count: 4 }).exec();
    const count5 = await TestModel.findOne({ count: 5 }).exec();
    expect(result.sliceInfo.firstCursor).toEqual(count4._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count5._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(true);
  });
  it('Find using before and after', async () => {
    const count3 = await TestModel.findOne({ count: 3 }).exec();
    const count4 = await TestModel.findOne({ count: 4 }).exec();
    const count5 = await TestModel.findOne({ count: 5 }).exec();
    const result = await findSlice({
      ...sliceQuery, after: count3._id.valueOf(),
    }, TestModel, altOption);
    expect(result.sliceInfo.firstCursor).toEqual(count4._id.valueOf());
    expect(result.sliceInfo.lastCursor).toEqual(count5._id.valueOf());
    expect(result.sliceInfo.hasNext).toEqual(false);
  });
});
