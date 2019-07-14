import { connectDb } from '~/utils/pg';

import app from '~/app';
import log from '~/utils/log';
import env from '~/utils/env';
import { connectMailgun } from '~/auth/mail';

(async () => {
  await connectDb(env.PG_URI);
  connectMailgun();
})();

log.info('run uexky at env:', env);
app.listen(env.PORT);
log.info(`Listening to http://localhost:${env.PORT} 🚀`);
