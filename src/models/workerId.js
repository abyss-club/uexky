// import Joi from 'joi';
import mongo from '~/utils/mongo';

// const workerIdSchema = Joi.object().keys({
//   count: Joi.number().integer().min(0),
// });

const WORKER_ID = 'workerid';

const WorkerIdModel = () => ({
  newWorkerId: async function newWorkerId() {
    const col = mongo.collection(WORKER_ID);
    const { value } = await col.findOneAndUpdate(
      {}, { $inc: { count: 1 } }, {
        projection: 'count', returnOriginal: false, upsert: true, w: 'majority', j: true, wtimeout: 1000,
      },
    );
    return value.count % 512;
  },
});

export default WorkerIdModel;
