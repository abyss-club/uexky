exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('user', {
    id: { type: 'serial', primaryKey: true },
    email: { type: 'text', notNull: true, unique: true },
    name: { type: 'text', unique: true },
    role: { type: 'text', notNull: true, default: '' },
    lastReadSystemNoti: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    lastReadRepliedNoti: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    lastReadQuotedNoti: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
  });
  pgm.createIndex('user', ['email']);
  pgm.createIndex('user', ['name']);

  pgm.createTable('thread', {
    id: { type: 'bigint', primaryKey: true },
    createdAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },

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
    createdAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
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
    createdAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    threadId: { type: 'bigint', notNull: true },
    userId: { type: 'integer', notNull: true },
    anonymousId: { type: 'bigint', notNull: true, unique: true },
  });
  pgm.createIndex('anonymous_id', ['threadId']);
  pgm.createIndex('anonymous_id', ['userId']);
  pgm.createIndex('anonymous_id', ['threadId', 'userId'], { unique: true });

  pgm.createTable('posts_quotes', {
    id: { type: 'serial', primaryKey: true },
    quoterId: { type: 'bigint', notNull: true, references: 'post(id)' },
    quotedId: { type: 'bigint', notNull: true, references: 'post(id)' },
  });
  pgm.createIndex('posts_quotes', ['quoterId', 'quotedId'], { unique: true });

  pgm.createTable('tag', {
    name: { type: 'text', primaryKey: true, notNull: true },
    isMain: { type: 'bool', notNull: true, default: false },
    createdAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
  });

  pgm.createTable('tags_main_tags', {
    id: { type: 'serial', primaryKey: true },
    createdAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    name: { type: 'text', notNull: true, references: 'tag(name)' },
    belongsTo: { type: 'text', notNull: true, references: 'tag(name)' },
  });
  pgm.createIndex('tags_main_tags', ['name', 'belongsTo'], { unique: true });

  pgm.createTable('threads_tags', {
    id: { type: 'serial', primaryKey: true },
    createdAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },
    tagName: { type: 'text', notNull: true, references: 'tag(name)' },
  });
  pgm.createIndex('threads_tags', ['threadId', 'tagName'], { unique: true });

  pgm.createTable('users_tags', {
    id: { type: 'serial', primaryKey: true },
    userId: { type: 'integer', notNull: true, references: 'public.user(id)' },
    tagName: { type: 'text', notNull: true, references: 'tag(name)' },
  });
  pgm.createIndex('users_tags', ['userId', 'tagName'], { unique: true });

  pgm.createTable('notification', {
    id: { type: 'text', primaryKey: true },
    createdAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updatedAt: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
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
  pgm.dropTable('anonymous_id');
  pgm.dropTable('post');
  pgm.dropTable('thread');
  pgm.dropTable('user');
};
