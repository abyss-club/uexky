exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('thread', {
    id: {
      type: 'string',
      primaryKey: true,
    },
    suid: { type: 'text', notNull: true, unique: true },
    anonymous: { type: 'bool', notNull: true },
    author: { type: 'text', notNull: true },
    userId: { type: 'text', notNull: true },
    mainTag: { type: 'text', notNull: true, references: '"tag" (name)' },
    subTags: { type: 'text[]', notNull: true },
    title: { type: 'text', notNull: true },
    locked: { type: 'bool', notNull: true },
    blocked: { type: 'bool', notNull: true },
    createdAt: { type: 'timestamp', notNull: true },
    updatedAt: { type: 'timestamp', notNull: true },
    content: { type: 'text', notNull: true },
  });
  // pgm.createIndex('suid', 'mainTag', 'subTags');
};

exports.down = (pgm) => {
  pgm.dropTable('thread');
};
