import Redis from 'ioredis';

import env from '~/utils/env';
import { ParamsError } from '~/utils/error';
import log from '~/utils/log';

async function createSubRateLimiter(config, forMutation, ip, email) {
  const key = email === ''
    ? `ratelimit:ip:${ip}`
    : `ratelimit:ip:${ip}:email:${email}`;
  const setKey = false;
  let redis = null;
  return {
    take: async function take(cost, isMutation = false) {
      if (!Number.isInteger(cost) || cost < 1) {
        throw new ParamsError('Limit must be greater than 0.');
      }
      if (forMutation !== isMutation) {
        return;
      }
      if (!redis) {
        log.info(`connect to redis ${env.REDIS_URI}`);
        redis = new Redis(env.REDIS_URI);
      }
      if (!setKey) {
        const {
          mutLimit, mutResetTime, queryLimit, queryResetTime,
        } = config.rateLimit;
        const limit = forMutation ? mutLimit : queryLimit;
        const expire = forMutation ? mutResetTime : queryResetTime;
        await redis.set(key, limit, 'EX', expire, 'NX');
      }
      const remaining = redis.decrby(key, cost);
      if (remaining < 0) {
        throw new Error('Rate limit exceeded.');
      }
    },
  };
}

async function createRateLimiter(config, ip, email = '') {
  let limiters = [];
  if (email !== '') {
    limiters = await Promise.all([
      createSubRateLimiter(config, false, ip, ''),
      createSubRateLimiter(config, true, ip, ''),
      createSubRateLimiter(config, false, ip, email),
      createSubRateLimiter(config, true, ip, email),
    ]);
  } else {
    limiters = [await createSubRateLimiter(config, false, ip, '')];
  }

  return {
    take: async function take(cost, isMutation = false) {
      await Promise.all(limiters.map(limiters.take(cost, isMutation)));
    },
  };
}

function createIdleRateLimiter() {
  return {
    take: function take() {},
  };
}

export default createRateLimiter;
export { createIdleRateLimiter };
