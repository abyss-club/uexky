import { connectDb } from '~/utils/pg';

import app from '~/app';
import log from '~/utils/log';
import env from '~/utils/env';

let pgPool;
(async () => {
  pgPool = await connectDb(env.PG_URI);
})();

log.info('run uexky at env:', env);
app.listen(env.PORT);
log.info(`Listening to http://localhost:${env.PORT} 🚀`);

(async () => {
  await pgPool.end();
})();
