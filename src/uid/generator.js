import mongoose from 'mongoose';

const randomSeq = () => Math.floor(Math.random() * 1024);
const timeZero = new Date('2018-03-01T00:00:00Z').getTime();
const timestamp = date => Math.floor((date.getTime() - timeZero) / 1000);

const WorkerIDSchema = new mongoose.Schema({
  count: Number,
}, { capped: 1 });
const workerExpireMilliSeconds = 1000 * 3600;

WorkerIDSchema.statics.newWorkerID = async function newWorkerID() {
  const { count } = await WorkerIDModel.findOneAndUpdate(
    {}, { $inc: { count: 1 } }, { new: true, upsert: 1 },
  ); // TODO: concernte(?) main, see mongodb docs.
  return count % 512;
};

const WorkerIDModel = mongoose.model('worker_id', WorkerIDSchema);

// Random Bits
const randomBits = () => Math.floor(Math.random() * 512);

const Generator = {
  workerID: '',
  expiredAt: 0,
  timestamp: 0,
  firstSeq: randomSeq(),
  seq: randomSeq(),
};

// run to next state, preparing for new id
Generator.run = async function run() {
  const now = new Date();
  if ((this.workerID === '') && (now.getTime() > this.expiredAt)) {
    this.workerID = await WorkerIDModel.newWorkerID();
    this.expiredAt = now + workerExpireMilliSeconds;
  }

  const nextSeq = (this.seq + 1) % 1024;
  const nowTs = timestamp(now);
  if (nowTs !== this.timestamp) {
    this.timestamp = nowTs;
    this.seq = nextSeq;
    this.firstSeq = nextSeq;
    return;
  }
  if (nextSeq !== this.firstSeq) {
    this.seq = nextSeq;
    return;
  }

  await setTimeout(1000 - now.getMilliseconds());
  this.timestamp += 1;
  this.seq = nextSeq;
  this.firstSeqInTs = nextSeq;
};

Generator.newID = async function newID() {
  await this.run();
  const rb = randomBits();
  return [
    this.timestamp.toString(16).padStart(8, '0'),
    (this.workerID * (2 ** 19)
      + this.seq * (2 ** 9) + rb).toString(16).padStart(7, '0'),
  ].join('');
};

export default Generator;
