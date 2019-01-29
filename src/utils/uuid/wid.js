import mongoose from 'mongoose';

const generator = {
  workerID: '',
  expiredAt: 0,
  ts: 0,
  firstSeqInTs: randomSeq(),
  seq: randomSeq(),
};

const randomSeq = () => Math.floor(Math.random() * 1024);

// Timestamp and Sequence
const timeZero = new Date('2018-03-01T00:00:00Z').getTime();
const timestamp = date => Math.floor((date.getTime() - timeZero) / 1000);
const sequenceNumber = async () => {
  const now = new Date();
  const nowTs = timestamp(now);

  const seq = (generator.seq + 1) % 1024;
  if (nowTs !== generator.ts) {
    generator.ts = nowTs;
    generator.seq = seq;
    generator.firstSeqInTs = seq;
    return { ts: nowTs, seq };
  }

  if (seq !== generator.firstSeqInTs) {
    generator.seq = seq;
    return { nowTs, seq };
  }

  await setTimeout(1000 - now.getMilliseconds);
  return sequenceNumber();
};

// Worker ID
const WorkerIDSchema = new mongoose.Schema({
  count: Number,
}, { capped: 1 });
const WorkerIDModel = mongoose.model('worker_id', WorkerIDSchema);

const expireMilliSeconds = 1000 * 3600;
const workerID = async () => {
  const now = new Date().getTime();
  if (generator.wid !== '' && generator.expiredAt > now) {
    return generator.wid;
  }

  const wid = await WorkerIDModel.findOneAndUpdate(
    {}, { $inc: { count: 1 } }, { new: true, upsert: 1 },
  );
  generate.wid = wid;
  generator.expiredAt = now + expireMilliSeconds;
  return generate.wid;
};


// Random Bits
const randomBits = () => Math.floor(Math.random() * 512);

const generate = async () => {
  const { seq, ts } = sequenceNumber();
  const wid = await workerID();
  const rb = randomBits();
  return ts * (2 ** 28) + wid * (2 ** 19) + seq * (2 ** 9) + rb;
};

export default generate;
