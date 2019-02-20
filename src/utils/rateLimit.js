import redis from '~/external/redis';
import { ParamsError } from '~/utils/error';

async function createSubRateLimiter(config, forMutation, ip, email) {
  const key = email === ''
    ? `ratelimit:ip:${ip}`
    : `ratelimit:ip:${ip}:email:${email}`;
  const set = false;
  return {
    take: async function take(cost, isMutation = false) {
      if (!Number.isInteger(cost) || cost < 1) {
        throw new ParamsError('Limit must be greater than 0.');
      }

      if (forMutation !== isMutation) return;
      if (!set) {
        const limitCfg = await config.getRateLimit();
        const limit = forMutation ? limitCfg.MutLimit : limitCfg.QueryLimit;
        const expire = forMutation
          ? limitCfg.MutResetTime : limitCfg.QueryResetTime;
        await redis.set(key, limit, 'EX', expire, 'NX');
      }
      const remaining = redis.decrby(key, cost);
      if (remaining < 0) {
        throw new Error('rate limit exceeded');
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
