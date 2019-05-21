exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.addColumns('workerId', {
    id: 'id',
  });
};

exports.down = (pgm) => {
  pgm.dropColumns('workerId', ['id']);
};
