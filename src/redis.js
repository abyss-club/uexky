import redis from 'redis';
import { promisify } from 'util';

const redisUrl = process.env.REDIS_URL;
const redisClient = (function createClient() {
  const client = redis.createClient(redisUrl);
  return {
    get: promisify(client.get).bind(client),
    set: promisify(client.set).bind(client),
    decrby: promisify(client.decrby).bind(client),
  };
}());

export default redisClient;
