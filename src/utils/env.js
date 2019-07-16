const env = {
  PG_URI: process.env.PG_URI || 'postgresql://localhost',
  REDIS_URI: process.env.REDIS_URI || 'redis://localhost:6379/0',

  MAILGUN_PRIVATE_KEY: process.env.MAILGUN_PRIVATE_KEY || 'private_key',
  MAILGUN_PUBLIC_KEY: process.env.MAILGUN_PUBLIC_KEY || 'public_key',
  MAILGUN_DOMAIN: process.env.MAILGUN_DOMAIN || 'mail.abyss.club',

  DOMAIN: process.env.DOMAIN,
  API_DOMAIN: process.env.API_DOMAIN,
  PROTO: process.env.PROTO || 'http',
  PORT: process.env.PORT || 5000,
  HOST: process.env.HOST || 'localhost',
};
export default env;
