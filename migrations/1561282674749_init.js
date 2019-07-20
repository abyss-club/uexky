exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable('user', {
    id: { type: 'integer', primaryKey: true, generated: { precedence: 'ALWAYS', increment: 1 } },
    email: { type: 'text', notNull: true, unique: true },
    name: { type: 'text', unique: true },
    role: { type: 'text', notNull: true, default: '' },
    last_read_system_noti: { type: 'bigint', notNull: true },
    last_read_replied_noti: { type: 'bigint', notNull: true },
    last_read_quoted_noti: { type: 'bigint', notNull: true },
  });
  pgm.createIndex('user', ['email']);
  pgm.createIndex('user', ['name']);

  pgm.createTable('thread', {
    id: { type: 'bigint', primaryKey: true },
    created_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updated_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },

    anonymous: { type: 'boolean', notNull: true },
    user_id: { type: 'integer', notNull: true, references: 'public.user(id)' },
    user_name: { type: 'varchar(16)', references: 'public.user(name)' },
    anonymous_id: { type: 'bigint' },

    title: { type: 'text', default: '' },
    content: { type: 'text', notNull: true },
    locked: { type: 'bool', notNull: true, default: false },
    blocked: { type: 'bool', notNull: true, default: false },
  });
  pgm.createIndex('thread', ['title']);
  pgm.createIndex('thread', ['anonymous']);
  pgm.createIndex('thread', ['user_id']);
  pgm.createIndex('thread', ['blocked']);

  pgm.createTable('post', {
    id: { type: 'bigint', primaryKey: true },
    created_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updated_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    thread_id: { type: 'bigint', notNull: true, references: 'thread(id)' },

    anonymous: { type: 'bool', notNull: true },
    user_id: { type: 'integer', notNull: true, references: 'public.user(id)' },
    user_name: { type: 'varchar(16)', references: 'public.user(name)' },
    anonymous_id: { type: 'bigint' },

    blocked: { type: 'bool', default: false },
    content: { type: 'text', notNull: true },
  });
  pgm.createIndex('post', ['thread_id']);
  pgm.createIndex('post', ['user_id']);

  pgm.createTable('anonymous_id', {
    id: { type: 'integer', primaryKey: true, generated: { precedence: 'ALWAYS', increment: 1 } },
    created_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updated_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    thread_id: { type: 'bigint', notNull: true },
    user_id: { type: 'integer', notNull: true },
    anonymous_id: { type: 'bigint', notNull: true, unique: true },
  });
  pgm.createIndex('anonymous_id', ['thread_id']);
  pgm.createIndex('anonymous_id', ['user_id']);
  pgm.createIndex('anonymous_id', ['thread_id', 'user_id'], { unique: true });

  pgm.createTable('posts_quotes', {
    id: { type: 'integer', primaryKey: true, generated: { precedence: 'ALWAYS', increment: 1 } },
    quoter_id: { type: 'bigint', notNull: true, references: 'post(id)' },
    quoted_id: { type: 'bigint', notNull: true, references: 'post(id)' },
  });
  pgm.createIndex('posts_quotes', ['quoter_id']);
  pgm.createIndex('posts_quotes', ['quoted_id']);
  pgm.createIndex('posts_quotes', ['quoter_id', 'quoted_id'], { unique: true });

  pgm.createTable('tag', {
    name: { type: 'text', primaryKey: true, notNull: true },
    is_main: { type: 'bool', notNull: true, default: false },
    created_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updated_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
  });
  pgm.createIndex('tag', ['is_main']);

  pgm.createTable('tags_main_tags', {
    id: { type: 'integer', primaryKey: true, generated: { precedence: 'ALWAYS', increment: 1 } },
    created_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updated_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    name: { type: 'text', notNull: true, references: 'tag(name)' },
    belongs_to: { type: 'text', notNull: true, references: 'tag(name)' },
  });
  pgm.createIndex('tag', ['name']);
  pgm.createIndex('tags_main_tags', ['name', 'belongs_to'], { unique: true });

  pgm.createTable('threads_tags', {
    id: { type: 'integer', primaryKey: true, generated: { precedence: 'ALWAYS', increment: 1 } },
    created_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updated_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    thread_id: { type: 'bigint', notNull: true, references: 'thread(id)' },
    tag_name: { type: 'text', notNull: true, references: 'tag(name)' },
  });
  pgm.createIndex('threads_tags', ['thread_id']);
  pgm.createIndex('threads_tags', ['thread_id', 'tag_name'], { unique: true });

  pgm.createTable('users_tags', {
    id: { type: 'integer', primaryKey: true, generated: { precedence: 'ALWAYS', increment: 1 } },
    user_id: { type: 'integer', notNull: true, references: 'public.user(id)' },
    tag_name: { type: 'text', notNull: true, references: 'tag(name)' },
  });
  pgm.createIndex('users_tags', ['user_id']);
  pgm.createIndex('users_tags', ['user_id', 'tag_name'], { unique: true });

  pgm.createTable('notification', {
    id: { type: 'bigint', primaryKey: true },
    key: { type: 'text', unique: true },
    created_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    updated_at: { type: 'timestamp with time zone', notNull: true, default: pgm.func('now()') },
    type: { type: 'text', notNull: true },
    send_to: { type: 'integer', references: 'public.user(id)' },
    send_to_group: { type: 'text' },
    content: { type: 'jsonb' },
  });
  pgm.createIndex('notification', ['type']);
  pgm.createIndex('notification', ['send_to']);
  pgm.createIndex('notification', ['send_to_group']);

  // pgm.createIndex('name');
  pgm.createTable('config', {
    id: { type: 'serial', primaryKey: true },
    rate_limit: { type: 'jsonb', notNull: true },
    rate_cost: { type: 'jsonb', notNull: true },
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
  pgm.dropTable('threads_tags');
  pgm.dropTable('tags_main_tags');
  pgm.dropTable('tag');
  pgm.dropTable('posts_quotes');
  pgm.dropTable('anonymous_id');
  pgm.dropTable('post');
  pgm.dropTable('thread');
  pgm.dropTable('user');
};
