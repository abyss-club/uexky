// import Joi from '@hapi/joi';
import { query } from '~/utils/pg';

const WorkerIdModel = () => ({
  newWorkerId: async function newWorkerId() {
    const results = await query('INSERT INTO workerid (id) VALUES(default) RETURNING id;');
    return results.rows[0].id % 512;
  },
});

export default WorkerIdModel;
