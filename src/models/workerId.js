import { query } from '~/utils/pg';

const WorkerIdModel = () => ({
  newWorkerId: async function newWorkerId() {
    const results = await query(
      "INSERT INTO counter (name) VALUES ('worker') ON CONFLICT (name)"
      + 'DO UPDATE SET count = counter.count + 1 RETURNING counter.count;',
    );
    return results.rows[0].count % 512;
  },
});

export default WorkerIdModel;
