// import Joi from 'joi';
import dbClient from '~/dbClient';

// const workerIdSchema = Joi.object().keys({
//   count: Joi.number().integer().min(0),
// });

const WORKER_ID = 'workerid';

const newWorkerId = async function newWorkerId() {
  const col = dbClient.collection(WORKER_ID);
  const { value } = await col.findOneAndUpdate(
    {}, { $inc: { count: 1 } }, {
      projection: 'count', returnOriginal: false, upsert: true, w: 'majority', j: true, wtimeout: 1000,
    },
  );
  return value.count % 512;
};

export { newWorkerId };
