const env = {
  PG_URI: process.env.PG_URI || 'postgresql://localhost',
  PG_DBNAME: process.env.PG_DBNAME || 'dev_uexky',
  REDIS_URI: process.env.REDIS_URI || 'redis://localhost:6379/0',

  MAINGUN_PRIVATE_KEY: process.env.MAINGUN_PRIVATE_KEY || 'private_key',
  MAINGUN_PUBLIC_KEY: process.env.MAINGUN_PUBLIC_KEY || 'public_key',
  MAINGUN_DOMAIN: process.env.MAINGUN_DOMAIN || 'mail.abyss.club',

  DOMAIN: process.env.DOMAIN,
  API_DOMAIN: process.env.API_DOMAIN,
  PROTO: process.env.PROTO || 'http',
  PORT: process.env.PORT || 5000,
};
export default env;
