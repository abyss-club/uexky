import { connectDb } from '~/utils/mongo';
import env from '~/utils/env';
import createIndexes from '~/models/indexes';


(async () => {
  const mongoClient = await connectDb(env.MONGODB_URI, env.MONGODB_DBNAME);
  await createIndexes();
  await mongoClient.close();
})();
