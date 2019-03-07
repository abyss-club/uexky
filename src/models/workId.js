import mongoose from 'mongoose';

const WorkerIDSchema = new mongoose.Schema(
  { count: Number },
  { writeConcern: { w: 'majority', j: true, wtimeout: 1000 } },
);

WorkerIDSchema.statics.newWorkerID = async function newWorkerID() {
  const { count } = await WorkerIDModel.findOneAndUpdate(
    {}, { $inc: { count: 1 } }, { new: true, upsert: 1 },
  ).exec();
  return count % 512;
};

const WorkerIDModel = mongoose.model('worker_id', WorkerIDSchema);

export default WorkerIDModel;
