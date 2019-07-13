import Redis from 'ioredis';

import env from '~/utils/env';
import log from '~/utils/log';

let redis = null;

function getRedis() {
  if (!redis) {
    log.info(`connect to redis ${env.REDIS_URI}`);
    redis = new Redis(env.REDIS_URI);
  }
  return redis;
}

export default getRedis;
