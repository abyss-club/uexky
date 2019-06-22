import WorkerIdModel from '~/models/workerId';

const randomSeq = () => Math.floor(Math.random() * 1024);
const timeZero = new Date('2018-03-01T00:00:00Z');
const timestamp = date => Math.floor((date.getTime() - timeZero.getTime()) / 1000);
const workerExpireMilliSeconds = 1000 * 3600;

// Random Bits
const randomBits = () => Math.floor(Math.random() * 512);

// async function newUID() -> BigInt
const newSuid = (function makeGenerator() {
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
      store.workerID = await WorkerIdModel().newWorkerId();
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

  return async () => {
    await run();
    const rb = randomBits();
    let id = BigInt(store.timestamp * (2 ** 28));
    id += BigInt(store.workerID * (2 ** 19));
    id += BigInt(store.seq * (2 ** 9));
    id += BigInt(rb);
    return id;
  };
}());

export default newSuid;
export { timeZero };
