-- DROP TABLE public.anonymous_id CASCADE;
-- DROP TABLE public.config CASCADE;
-- DROP TABLE public.counter CASCADE;
-- DROP TABLE public.pgmigrations CASCADE;

ALTER TABLE public.post ADD COLUMN quoted_ids bigint[];
CREATE INDEX post_quote_post_index ON public.post USING gin (quoted_ids);

-- DROP TABLE public.posts_quotes CASCADE;

ALTER TABLE public.tag ADD COLUMN tag_type text;
ALTER TABLE public.tag ALTER COLUMN is_main DROP NOT NULL;
-- DROP INDEX tag_is_main_index;
-- DROP ALTER TABLE public.tag DROP COLUMN is_main;

-- DROP TABLE public.tags_main_tags CASCADE;

ALTER TABLE public.thread ADD COLUMN tags text[];
CREATE INDEX thread_tags_index ON public.thread USING gin (tags);

-- DROP TABLE public.threads_tags CASCADE;

ALTER TABLE public."user" ADD COLUMN tags text[];
ALTER TABLE public."user" ADD COLUMN last_read_noti bigint DEFAULT 0 NOT NULL;
CREATE INDEX user_tags_index ON public."user" USING gin (tags);

-- DROP TABLE users_tags CASCADE;
