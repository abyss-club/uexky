CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    return NEW;
end;
$$ language 'plpgsql';


-- user


CREATE TABLE public."user" (
    id bigint PRIMARY KEY,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    email text UNIQUE,
    name text,
    role text NOT NULL,
    last_read_noti bigint NOT NULL,
    tags text[]
);

CREATE TRIGGER user_updated_at
    before update on public."user"
    for each row
    execute procedure update_updated_at();


-- thread


CREATE TABLE public.thread (
    id bigint PRIMARY KEY,
    last_post_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL,
    anonymous boolean NOT NULL,
    guest boolean NOT NULL,
    author character varying(16) NOT NULL,
    title text DEFAULT ''::text,
    content text NOT NULL,
    locked boolean DEFAULT false NOT NULL,
    blocked boolean DEFAULT false NOT NULL,
    tags text[]
);

CREATE TRIGGER thread_updated_at
    before update on public.thread
    for each row
    execute procedure update_updated_at();

CREATE INDEX thread_last_post_id_index ON public.thread USING btree (last_post_id);
CREATE INDEX thread_user_id_index ON public.thread USING btree (user_id);
CREATE INDEX thread_tags_index ON public.thread USING gin (tags);


-- post


CREATE TABLE public.post (
    id bigint PRIMARY KEY,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    thread_id bigint NOT NULL,
    user_id bigint NOT NULL,
    anonymous boolean NOT NULL,
    guest boolean NOT NULL,
    author character varying(16) NOT NULL,
    blocked boolean NOT NULL DEFAULT false,
    content text NOT NULL,
    quoted_ids bigint[]
);

CREATE TRIGGER post_updated_at
    before update on public.post
    for each row
    execute procedure update_updated_at();

CREATE INDEX post_thread_id_index ON public.post USING btree (thread_id);
CREATE INDEX post_user_id_index ON public.post USING btree (user_id);
CREATE INDEX post_quote_post_index ON public.post USING gin (quoted_ids);


-- tag


CREATE TABLE tag (
    name text PRIMARY KEY,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    type text NOT NULL
);

CREATE TRIGGER tag_updated_at
    before update on public.tag
    for each row
    execute procedure update_updated_at();


-- notification


CREATE TABLE notification (
    key text PRIMARY KEY,
    sort_key bigint NOT NULL UNIQUE,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    type text NOT NULL,
    receivers text[],
    content jsonb
);

CREATE TRIGGER notification_updated_at
    before update on public.notification
    for each row
    execute procedure update_updated_at();

CREATE INDEX notification_reveivers_index ON public.notification USING gin (receivers);
