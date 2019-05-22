exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('workerId', {
    id: {
      type: 'integer',
      primaryKey: true,
      generated: {
        precedence: 'BY DEFAULT',
        increment: 1,
      },
    },
  });
};

exports.down = (pgm) => {
  pgm.dropTable('workerId');
};
