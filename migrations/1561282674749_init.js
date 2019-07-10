exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('user', {
    id: { type: 'serial', primaryKey: true },
    email: { type: 'text', notNull: true, unique: true },
    name: { type: 'text', unique: true },
    role: { type: 'text', notNull: true, default: '' },
    lastReadSystemNoti: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    lastReadRepliedNoti: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    lastReadQuotedNoti: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
  });
  pgm.createIndex('user', ['email']);
  pgm.createIndex('user', ['name']);

  pgm.createTable('thread', {
    id: { type: 'bigint', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },

    anonymous: { type: 'boolean', notNull: true },
    userId: { type: 'integer', notNull: true, references: 'public.user(id)' },
    userName: { type: 'varchar(16)', references: 'public.user(name)' },
    anonymousId: { type: 'bigint' },

    title: { type: 'text', default: '' },
    content: { type: 'text', notNull: true },
    locked: { type: 'bool', notNull: true, default: false },
    blocked: { type: 'bool', notNull: true, default: false },
  });
  pgm.createIndex('thread', 'title');
  pgm.createIndex('thread', 'anonymous');
  pgm.createIndex('thread', 'userId');
  pgm.createIndex('thread', 'blocked');

  pgm.createTable('post', {
    id: { type: 'bigint', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },

    anonymous: { type: 'bool', notNull: true },
    userId: { type: 'integer', notNull: true, references: 'public.user(id)' },
    userName: { type: 'varchar(16)', references: 'public.user(name)' },
    anonymousId: { type: 'bigint' },

    blocked: { type: 'bool', default: false },
    content: { type: 'text', notNull: true },
  });
  pgm.createIndex('post', ['threadId', 'userId', 'anonymous']);

  pgm.createTable('anonymous_id', {
    id: { type: 'serial', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },
    userId: { type: 'integer', notNull: true, references: 'public.user(id)' },
    anonymousId: { type: 'bigint', notNull: true, unique: true },
  });
  pgm.addConstraint('thread', 'thread_anonymous_id', {
    foreignKeys: {
      columns: 'anonymousId',
      references: 'anonymous_id("anonymousId")',
    },
  });
  pgm.addConstraint('post', 'post_anonymous_id', {
    foreignKeys: {
      columns: 'anonymousId',
      references: 'anonymous_id("anonymousId")',
    },
  });

  pgm.createTable('posts_quotes', {
    id: { type: 'serial', primaryKey: true },
    quoterId: { type: 'bigint', notNull: true, references: 'post(id)' },
    quotedId: { type: 'bigint', notNull: true, references: 'post(id)' },
  });

  pgm.createTable('tag', {
    name: { type: 'text', primaryKey: true, notNull: true },
    isMain: { type: 'bool', notNull: true },
    createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
  });

  pgm.createTable('tags_main_tags', {
    id: { type: 'serial', primaryKey: true },
    name: { type: 'text', notNull: true, references: 'tag(name)' },
    belongsTo: { type: 'text', notNull: true, references: 'tag(name)' },
  });

  pgm.createTable('threads_tags', {
    id: { type: 'serial', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },
    tagName: { type: 'text', notNull: true, references: 'tag(name)' },
  });

  pgm.createTable('users_tags', {
    id: { type: 'serial', primaryKey: true },
    userId: { type: 'integer', notNull: true, references: 'public.user(id)' },
    tagName: { type: 'text', notNull: true, references: 'tag(name)' },
  });
  pgm.createIndex('users_tags', ['userId', 'tagName'], { unique: true });

  pgm.createTable('notification', {
    id: { type: 'serial', primaryKey: true },
    createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
    type: { type: 'text', notNull: true },
    sendTo: { type: 'integer', references: 'public.user(id)' },
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
  pgm.dropTable('counter');
  pgm.dropTable('config');
  pgm.dropTable('notification');
  pgm.dropTable('users_tags');
  pgm.dropTable('tags_main_tags');
  pgm.dropConstraint('thread_main_tag');
  pgm.dropTable('threads_sub_tags');
  pgm.dropTable('tag');
  pgm.dropTable('posts_quotes');
  pgm.dropConstraint('thread', 'thread_anonymous_id');
  pgm.dropConstraint('post', 'post_anonymous_id');
  pgm.dropTable('anonymous_id');
  pgm.dropTable('post');
  pgm.dropTable('thread');
  pgm.dropTable('user');
};
