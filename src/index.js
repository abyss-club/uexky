import { connectDb } from '~/utils/mongo';

import app from '~/app';
import log from '~/utils/log';
import env from '~/utils/env';

let mongoClient;
(async () => {
  mongoClient = await connectDb(env.MONGODB_URI, env.MONGODB_DBNAME);
})();

log.info('run uexky at env:', env);
app.listen(env.PORT);
log.info(`Listening to http://localhost:${env.PORT} 🚀`);

(async () => {
  await mongoClient.close();
})();
