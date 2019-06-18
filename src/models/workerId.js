import { query } from '~/utils/pg';

/* CREATE table counter(
 *     name varchar(32) PRIMARY KEY,
 *     count integer DEFAULT 0
 * );
*/

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
