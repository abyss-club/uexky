import mongoose from 'mongoose';

import env from '~/utils/env';
import log from '~/utils/log';
import createIndexes from '~/models/indexes';

const result = (async () => {
  await mongoose.connect(env.MONGODB_URI, {
    useNewUrlParser: true,
    useCreateIndex: true,
    useFindAndModify: false,
    autoIndex: false,
  });
  await createIndexes();
  await mongoose.disconnect();
})();

log.info(`${result}`);
