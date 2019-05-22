exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('config', {
    id: {
      type: 'integer',
      primaryKey: true,
      generated: {
        precedence: 'BY DEFAULT',
        increment: 1,
      },
    },
    rateLimit: {
      type: 'json',
      notNull: true,
    },
    rateCost: {
      type: 'json',
      notNull: true,
    },
  });
};

exports.down = (pgm) => {
  pgm.dropTable('config');
};
