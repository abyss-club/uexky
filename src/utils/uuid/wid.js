import mongoose from 'mongoose';

// Worker id

const WorkerIDSchema = new mongoose.Schema({
  count: Number,
}, { capped: 1 });
const WorkerIDModel = mongoose.model('worker_id', WorkerIDSchema);

const pidStore = {
  wid: '',
  expiredAt: 0,
};

function expireSeconds() {
  return 3600 + Math.floor(Math.random() * 3600);
}

async function getWorkerID() {
  const now = new Date().getTime();
  if (pidStore.wid !== '' && pidStore.expiredAt < now) {
    return pidStore.wid;
  }

  pidStore.expiredAt = now + 1000 * expireSeconds();
  let wid = '';
  try {
    wid = await WorkerIDModel.findOneAndUpdate(
      {}, { $inc: { count: 1 } }, { new: true, upsert: 1 },
    );
  } catch (e) {
    throw e;
  }
  pidStore.wid = wid;
  return wid;
}

export default getWorkerID;
