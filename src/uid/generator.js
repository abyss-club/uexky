import mongoose from 'mongoose';

const randomSeq = () => Math.floor(Math.random() * 1024);
const timeZero = new Date('2018-03-01T00:00:00Z').getTime();
const timestamp = date => Math.floor((date.getTime() - timeZero) / 1000);

const WorkerIDSchema = new mongoose.Schema(
  { count: Number },
  { capped: 1, writeConcern: { w: 'majority', j: true, wtimeout: 1000 } },
);
const workerExpireMilliSeconds = 1000 * 3600;

WorkerIDSchema.statics.newWorkerID = async function newWorkerID() {
  const { count } = await WorkerIDModel.findOneAndUpdate(
    {}, { $inc: { count: 1 } }, { new: true, upsert: 1 },
  );
  return count % 512;
};

const WorkerIDModel = mongoose.model('worker_id', WorkerIDSchema);

// Random Bits
const randomBits = () => Math.floor(Math.random() * 512);

const generator = (function makeGenerator() {
  const store = {
    workerID: '',
    expiredAt: 0,
    timestamp: 0,
    firstSeq: randomSeq(),
    seq: randomSeq(),
  };
  const run = async () => {
    const now = new Date();
    if ((store.workerID === '') && (now.getTime() > store.expiredAt)) {
      store.workerID = await WorkerIDModel.newWorkerID();
      store.expiredAt = now + workerExpireMilliSeconds;
    }

    const nextSeq = (store.seq + 1) % 1024;
    const nowTs = timestamp(now);
    if (nowTs !== store.timestamp) {
      store.timestamp = nowTs;
      store.seq = nextSeq;
      store.firstSeq = nextSeq;
      return;
    }
    if (nextSeq !== store.firstSeq) {
      store.seq = nextSeq;
      return;
    }

    await setTimeout(1000 - now.getMilliseconds());
    store.timestamp += 1;
    store.seq = nextSeq;
    store.firstSeqInTs = nextSeq;
  };

  const newID = async () => {
    await run();
    const rb = randomBits();
    return [
      store.timestamp.toString(16).padStart(8, '0'),
      (store.workerID * (2 ** 19)
      + store.seq * (2 ** 9) + rb).toString(16).padStart(7, '0'),
    ].join('');
  };
  return { newID };
}());

export default generator;
export { WorkerIDModel };
