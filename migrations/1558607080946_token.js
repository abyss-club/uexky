exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('token', {
    id: {
      type: 'integer',
      primaryKey: true,
      generated: {
        precedence: 'BY DEFAULT',
        increment: 1,
      },
    },
    email: { type: 'text', unique: true, notNull: true },
    authToken: { type: 'text', notNull: true },
  });
  pgm.createIndex('token', 'email', 'authToken');
};

exports.down = (pgm) => {
  pgm.dropTable('token');
};
