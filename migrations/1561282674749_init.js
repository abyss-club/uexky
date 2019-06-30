exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('user', {
    id: { type: 'serial', primaryKey: true },
    email: { type: 'text', notNull: true, unique: true },
    name: { type: 'text', unique: true },
    role: { type: 'text' },
    lastReadSystemNoti: { type: 'timestamp', default: 'now()' },
    lastReadRepliedNoti: { type: 'timestamp', default: 'now()' },
    lastReadQuotedNoti: { type: 'timestamp', default: 'now()' },
  });

  pgm.createTable('thread', {
    id: { type: 'bigint', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true },
    updatedAt: { type: 'timestamp', notNull: true },

    anonymous: { type: 'boolean', notNull: true },
    userId: { type: 'integer', notNull: true, references: 'user(id)' },
    userName: { type: 'varchar(16)', references: 'user(name)' },
    anonymousId: { type: 'bigint', references: 'anonymous_id(anonymous_id)' },

    title: { type: 'text', notNull: true },
    locked: { type: 'bool', notNull: true },
    blocked: { type: 'bool', notNull: true },
    content: { type: 'text', notNull: true },
  });
  pgm.createIndex('title', 'anonymous', 'userId', 'blocked');

  pgm.createTable('post', {
    id: { type: 'bigint', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true },
    updatedAt: { type: 'timestamp', notNull: true },
    threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },

    anonymous: { type: 'bool', notNull: true },
    userId: { type: 'integer', notNull: true, references: 'user(id)' },
    userName: { type: 'varchar(16)', references: 'user(name)' },
    anonymousId: { type: 'bigint', references: 'anonymous_id(anonymousId)' },

    locked: { type: 'bool', default: false },
    blocked: { type: 'bool', default: false },
    content: { type: 'text', notNull: true },
  });
  pgm.createIndex('threadId', 'userId', 'anonymous');

  pgm.createTable('anonymous_id', {
    id: { type: 'serial', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true },
    updatedAt: { type: 'timestamp', notNull: true },
    threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },
    userId: { type: 'integer', notNull: true, references: 'user(id)' },
    anonymousId: { type: 'bigint', notNull: true },
  });

  pgm.createTable('tag', {
    name: { type: 'text', primaryKey: true, notNull: true },
    isMain: { type: 'bool', notNull: true },
    createdAt: { type: 'timestamp', notNull: true },
    updatedAt: { type: 'timestamp', notNull: true },
  });

  pgm.createTable('threads_tags', {
    id: { type: 'serial', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true },
    updatedAt: { type: 'timestamp', notNull: true },
    threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },
    tagName: { type: 'text', notNull: true, references: 'tag(name)' },
  });

  pgm.createTable('users_tags', {
    id: { type: 'serial', primaryKey: true },
    userId: { type: 'integer', notNull: true, references: 'user(id)' },
    tagName: { type: 'text', notNull: true, references: 'tag(name)' },
  });

  pgm.createTable('managers_tags', {
    id: { type: 'serial', primaryKey: true },
    userId: { type: 'integer', notNull: true, references: 'user(id)' },
    tagName: { type: 'text', notNull: true, references: 'tag(name)' },
  });

  pgm.createTable('notification', {
    id: { type: 'serial', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true },
    type: { type: 'text', notNull: true },
    sendTo: { type: 'integer', references: 'user(id)' },
    sendToGroup: { type: 'text' },
    content: { type: 'jsonb' },
  });

  // pgm.createIndex('name');
  pgm.createTable('config', {
    id: { type: 'serial', primaryKey: true },
    rateLimit: { type: 'jsonb', notNull: true },
    rateCost: { type: 'jsonb', notNull: true },
  });
  pgm.createTable('counter', {
    name: { type: 'varchar(32)', primaryKey: true },
    count: { type: 'integer', default: 0 },
  });
};

exports.down = (pgm) => {
  pgm.dropTable('user');
  pgm.dropTable('thread');
  pgm.dropTable('post');
  pgm.dropTable('anonymous_id');
  pgm.dropTable('tag');
  pgm.dropTable('threads_tags');
  pgm.dropTable('users_tags');
  pgm.dropTable('managers_tags');
  pgm.dropTable('notification');
  pgm.dropTable('config');
  pgm.dropTable('counter');
};
