import mongoose from 'mongoose';

import env from '~/utils/env';
import log from '~/utils/log';
import createIndexes from '~/models/indexes';

mongoose.connect(env.MONGODB_URI, {
  useNewUrlParser: true,
  useCreateIndex: true,
  useFindAndModify: false,
});

createIndexes().then(() => {
  log.info('created indexes!');
  process.exit();
}).catch((err) => {
  log.error(err);
  process.exit(1);
});
