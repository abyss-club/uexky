import mongoose from 'mongoose';

import app from '~/app';
import log from '~/utils/log';
import env from '~/utils/env';

mongoose.connect(env.MONGODB_URI, {
  useNewUrlParser: true,
  useCreateIndex: true,
  useFindAndModify: false,
});

app.listen(env.PORT);
log.info(`Listening to http://localhost:${env.PORT} 🚀`);
