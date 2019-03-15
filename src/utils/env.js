const env = {
  MONGODB_URI: process.env.MONGODB_URI || 'mongodb://localhost:27017,localhost:27018,localhost:27019/dev_uexky?replicaSet=rs',
  MONGODB_DBNAME: process.env.MONGODB_DBNAME || 'dev_uexky',
  REDIS_URI: process.env.REDIS_URI || 'redis://localhost:6379/0',

  MAINGUN_PRIVATE_KEY: process.env.MAINGUN_PRIVATE_KEY,
  MAINGUN_DOMAIN: process.env.MAINGUN_DOMAIN,

  DOMAIN: process.env.DOMAIN,
  API_DOMAIN: process.env.API_DOMAIN,
  PROTO: process.env.PROTO || 'http',
  PORT: process.env.PORT || 5000,
};
export default env;
