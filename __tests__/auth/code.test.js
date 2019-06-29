import Code from '~/auth/code';
import getRedis from '~/utils/redis';

import mockMailgun from '../__utils__/mailgun';

afterAll(async () => {
  const redis = getRedis();
  redis.flushall();
});

describe('Testing auth', () => {
  const mockEmail = 'test@example.com';
  it('add user', async () => {
    const mailgun = mockMailgun();
    const code = await Code.addToAuth(mockEmail);
    const email = await Code.getEmailByCode(code);
    expect(email).toEqual(mockEmail);
    expect(mailgun.mail.to).toEqual(mockEmail);
    expect(mailgun.mail.text).toMatch(code);
  });
});
