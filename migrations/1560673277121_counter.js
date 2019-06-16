exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('counter', {
    name: {
      type: 'varchar(32)',
      primaryKey: true,
    },
    count: {
      type: 'integer',
      default: 0,
    },
  });
};

exports.down = (pgm) => {
  pgm.dropTable('counter');
};
