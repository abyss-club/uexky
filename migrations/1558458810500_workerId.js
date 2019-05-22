exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('workerId', {
    id: 'id',
  });
};

exports.down = (pgm) => {
  pgm.dropTable('workerId');
};
