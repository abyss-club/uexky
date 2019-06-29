import Token from '~/auth/token';
import getRedis from '~/utils/redis';

afterAll(async () => {
  const redis = getRedis();
  await redis.flushall();
});

describe('Testing token', () => {
  const mockEmail = 'test@example.com';
  it('validate token by email', async () => {
    const token = await Token.genNewToken(mockEmail);
    const email = await Token.getEmailByToken(token);
    expect(email).toEqual(mockEmail);
  });
});
