import winston from 'winston';

const logger = winston.createLogger({
  level: 'debug',
  format: winston.format.combine(
    winston.format.errors({ stack: true }),
    winston.format.timestamp({ alias: 'time' }),
    winston.format.json(),
    winston.format.prettyPrint(),
  ),
  transports: [new winston.transports.Console({
    stderrLevels: ['error'],
  })],
});

const log = {
  info: (message, data) => { logger.info(message, data); },
  warn: (message, data) => { logger.warn(message, data); },
  error: (message, data) => {
    logger.error(message, data);
    if (message instanceof Error) {
      /* eslint-disable no-console */
      console.log(message.stackSourcegraph);
      /* eslint-enable */
    }
  },
  debug: (message, data) => {
    logger.log({ ...data, level: 'debug', message });
  },
};

export default log;
