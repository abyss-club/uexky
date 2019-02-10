import redis from '~/redis';

async function createSubRateLimiter(config, forMutation, ip, email) {
  const key = email === ''
    ? `ratelimit:ip:${ip}`
    : `ratelimit:ip:${ip}:email:${email}`;
  const set = false;
  return {
    take: async function take(cost, isMutation) {
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
  const limiters = [createSubRateLimiter(config, false, ip, '')];
  if (email !== '') {
    limiters.push(createSubRateLimiter(config, true, ip, ''));
    limiters.push(createSubRateLimiter(config, false, ip, email));
    limiters.push(createSubRateLimiter(config, true, ip, email));
  }
  return {
    take: async function take(cost, isMutation) {
      await Promise.all(limiters.map(limiters.take(cost, isMutation)));
    },
  };
}


export default createRateLimiter;
