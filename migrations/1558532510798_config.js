exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('config', {
    id: 'id',
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
