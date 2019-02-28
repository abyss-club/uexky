import redis from 'redis';
import { promisify } from 'util';

import env from '~/utils/env';

const redisClient = (function createClient() {
  const client = redis.createClient(env.REDIS_URI);
  return {
    get: promisify(client.get).bind(client),
    set: promisify(client.set).bind(client),
    decrby: promisify(client.decrby).bind(client),
    quit: () => { client.quit(); },
  };
}());

export default redisClient;
