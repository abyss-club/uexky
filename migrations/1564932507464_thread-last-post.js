exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.addColumns('thread', { last_post_id: { type: 'bigint', notNull: true, default: 0 } });
  pgm.createIndex('thread', ['last_post_id']);
};

exports.down = (pgm) => {
  pgm.dropIndex('thread', ['last_post_id']);
  pgm.dropColumns('thread', ['last_post_id']);
};
