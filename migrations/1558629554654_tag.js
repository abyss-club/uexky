exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('tag', {
    name: {
      type: 'text',
      primaryKey: true,
      unique: true,
      notNull: true,
    },
    updatedAt: {
      type: 'timestamp',
    },
  });
  // pgm.createIndex('name');
};

exports.down = (pgm) => {
  pgm.dropTable('tag');
};
