exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('workerid', {
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
  pgm.dropTable('workerid');
};
