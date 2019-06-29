import Token from '~/auth/token';
import getRedis from '~/utils/redis';

afterAll(() => {
  const redis = getRedis();
  redis.flushall();
  redis.disconnect();
});

describe('Testing token', () => {
  const mockEmail = 'test@example.com';
  it('validate token by email', async () => {
    const token = await Token.genNewToken(mockEmail);
    const email = await Token.getEmailByToken(token);
    expect(email).toEqual(mockEmail);
  });
});
