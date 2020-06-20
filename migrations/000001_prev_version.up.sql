--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2 (Debian 12.2-2.pgdg100+1)
-- Dumped by pg_dump version 12.2 (Debian 12.2-2.pgdg100+1)

-- SET statement_timeout = 0;
-- SET lock_timeout = 0;
-- SET idle_in_transaction_session_timeout = 0;
-- SET client_encoding = 'UTF8';
-- SET standard_conforming_strings = on;
-- SELECT pg_catalog.set_config('search_path', '', false);
-- SET check_function_bodies = false;
-- SET xmloption = content;
-- SET client_min_messages = warning;
-- SET row_security = off;
-- 
-- SET default_tablespace = '';
-- 
-- SET default_table_access_method = heap;

--
-- Name: anonymous_id; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.anonymous_id (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    thread_id bigint NOT NULL,
    user_id integer NOT NULL,
    anonymous_id bigint NOT NULL
);


ALTER TABLE public.anonymous_id OWNER TO postgres;

--
-- Name: anonymous_id_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.anonymous_id ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.anonymous_id_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.config (
    id integer NOT NULL,
    rate_limit jsonb NOT NULL,
    rate_cost jsonb NOT NULL
);


ALTER TABLE public.config OWNER TO postgres;

--
-- Name: config_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.config_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.config_id_seq OWNER TO postgres;

--
-- Name: config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.config_id_seq OWNED BY public.config.id;


--
-- Name: counter; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.counter (
    name character varying(32) NOT NULL,
    count integer DEFAULT 0
);


ALTER TABLE public.counter OWNER TO postgres;

--
-- Name: notification; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notification (
    id integer NOT NULL,
    key text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    type text NOT NULL,
    send_to integer,
    send_to_group text,
    content jsonb
);


ALTER TABLE public.notification OWNER TO postgres;

--
-- Name: notification_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.notification ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.notification_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: pgmigrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pgmigrations (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    run_on timestamp without time zone NOT NULL
);


ALTER TABLE public.pgmigrations OWNER TO postgres;

--
-- Name: pgmigrations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pgmigrations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pgmigrations_id_seq OWNER TO postgres;

--
-- Name: pgmigrations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pgmigrations_id_seq OWNED BY public.pgmigrations.id;


--
-- Name: post; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.post (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    thread_id bigint NOT NULL,
    anonymous boolean NOT NULL,
    user_id integer NOT NULL,
    user_name character varying(16),
    anonymous_id bigint,
    blocked boolean DEFAULT false,
    content text NOT NULL
);


ALTER TABLE public.post OWNER TO postgres;

--
-- Name: posts_quotes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.posts_quotes (
    id integer NOT NULL,
    quoter_id bigint NOT NULL,
    quoted_id bigint NOT NULL
);


ALTER TABLE public.posts_quotes OWNER TO postgres;

--
-- Name: posts_quotes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.posts_quotes ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.posts_quotes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: tag; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tag (
    name text NOT NULL,
    is_main boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.tag OWNER TO postgres;

--
-- Name: tags_main_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tags_main_tags (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    name text NOT NULL,
    belongs_to text NOT NULL
);


ALTER TABLE public.tags_main_tags OWNER TO postgres;

--
-- Name: tags_main_tags_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.tags_main_tags ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.tags_main_tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: thread; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.thread (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    anonymous boolean NOT NULL,
    user_id integer NOT NULL,
    user_name character varying(16),
    anonymous_id bigint,
    title text DEFAULT ''::text,
    content text NOT NULL,
    locked boolean DEFAULT false NOT NULL,
    blocked boolean DEFAULT false NOT NULL,
    last_post_id bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public.thread OWNER TO postgres;

--
-- Name: threads_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.threads_tags (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    thread_id bigint NOT NULL,
    tag_name text NOT NULL
);


ALTER TABLE public.threads_tags OWNER TO postgres;

--
-- Name: threads_tags_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.threads_tags ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.threads_tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    email text NOT NULL,
    name text,
    role text DEFAULT ''::text NOT NULL,
    last_read_system_noti integer DEFAULT 0 NOT NULL,
    last_read_replied_noti integer DEFAULT 0 NOT NULL,
    last_read_quoted_noti integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."user" ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_tags (
    id integer NOT NULL,
    user_id integer NOT NULL,
    tag_name text NOT NULL
);


ALTER TABLE public.users_tags OWNER TO postgres;

--
-- Name: users_tags_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.users_tags ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: config id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config ALTER COLUMN id SET DEFAULT nextval('public.config_id_seq'::regclass);


--
-- Name: pgmigrations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pgmigrations ALTER COLUMN id SET DEFAULT nextval('public.pgmigrations_id_seq'::regclass);


--
-- Name: anonymous_id anonymous_id_anonymous_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.anonymous_id
    ADD CONSTRAINT anonymous_id_anonymous_id_key UNIQUE (anonymous_id);


--
-- Name: anonymous_id anonymous_id_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.anonymous_id
    ADD CONSTRAINT anonymous_id_pkey PRIMARY KEY (id);


--
-- Name: config config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_pkey PRIMARY KEY (id);


--
-- Name: counter counter_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.counter
    ADD CONSTRAINT counter_pkey PRIMARY KEY (name);


--
-- Name: notification notification_key_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification
    ADD CONSTRAINT notification_key_key UNIQUE (key);


--
-- Name: notification notification_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification
    ADD CONSTRAINT notification_pkey PRIMARY KEY (id);


--
-- Name: pgmigrations pgmigrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pgmigrations
    ADD CONSTRAINT pgmigrations_pkey PRIMARY KEY (id);


--
-- Name: post post_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post
    ADD CONSTRAINT post_pkey PRIMARY KEY (id);


--
-- Name: posts_quotes posts_quotes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts_quotes
    ADD CONSTRAINT posts_quotes_pkey PRIMARY KEY (id);


--
-- Name: tag tag_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tag
    ADD CONSTRAINT tag_pkey PRIMARY KEY (name);


--
-- Name: tags_main_tags tags_main_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags_main_tags
    ADD CONSTRAINT tags_main_tags_pkey PRIMARY KEY (id);


--
-- Name: thread thread_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thread
    ADD CONSTRAINT thread_pkey PRIMARY KEY (id);


--
-- Name: threads_tags threads_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads_tags
    ADD CONSTRAINT threads_tags_pkey PRIMARY KEY (id);


--
-- Name: user user_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_email_key UNIQUE (email);


--
-- Name: user user_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_name_key UNIQUE (name);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: users_tags users_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_tags
    ADD CONSTRAINT users_tags_pkey PRIMARY KEY (id);


--
-- Name: anonymous_id_thread_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX anonymous_id_thread_id_index ON public.anonymous_id USING btree (thread_id);


--
-- Name: anonymous_id_thread_id_user_id_unique_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX anonymous_id_thread_id_user_id_unique_index ON public.anonymous_id USING btree (thread_id, user_id);


--
-- Name: anonymous_id_user_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX anonymous_id_user_id_index ON public.anonymous_id USING btree (user_id);


--
-- Name: notification_created_at_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX notification_created_at_index ON public.notification USING btree (created_at);


--
-- Name: notification_send_to_group_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX notification_send_to_group_index ON public.notification USING btree (send_to_group);


--
-- Name: notification_send_to_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX notification_send_to_index ON public.notification USING btree (send_to);


--
-- Name: notification_type_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX notification_type_index ON public.notification USING btree (type);


--
-- Name: post_thread_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_thread_id_index ON public.post USING btree (thread_id);


--
-- Name: post_user_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_user_id_index ON public.post USING btree (user_id);


--
-- Name: posts_quotes_quoted_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX posts_quotes_quoted_id_index ON public.posts_quotes USING btree (quoted_id);


--
-- Name: posts_quotes_quoter_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX posts_quotes_quoter_id_index ON public.posts_quotes USING btree (quoter_id);


--
-- Name: posts_quotes_quoter_id_quoted_id_unique_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX posts_quotes_quoter_id_quoted_id_unique_index ON public.posts_quotes USING btree (quoter_id, quoted_id);


--
-- Name: tag_is_main_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX tag_is_main_index ON public.tag USING btree (is_main);


--
-- Name: tags_main_tags_name_belongs_to_unique_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX tags_main_tags_name_belongs_to_unique_index ON public.tags_main_tags USING btree (name, belongs_to);


--
-- Name: tags_main_tags_name_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX tags_main_tags_name_index ON public.tags_main_tags USING btree (name);


--
-- Name: thread_anonymous_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX thread_anonymous_index ON public.thread USING btree (anonymous);


--
-- Name: thread_blocked_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX thread_blocked_index ON public.thread USING btree (blocked);


--
-- Name: thread_last_post_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX thread_last_post_id_index ON public.thread USING btree (last_post_id);


--
-- Name: thread_title_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX thread_title_index ON public.thread USING btree (title);


--
-- Name: thread_user_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX thread_user_id_index ON public.thread USING btree (user_id);


--
-- Name: threads_tags_thread_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX threads_tags_thread_id_index ON public.threads_tags USING btree (thread_id);


--
-- Name: threads_tags_thread_id_tag_name_unique_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX threads_tags_thread_id_tag_name_unique_index ON public.threads_tags USING btree (thread_id, tag_name);


--
-- Name: user_email_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX user_email_index ON public."user" USING btree (email);


--
-- Name: user_name_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX user_name_index ON public."user" USING btree (name);


--
-- Name: users_tags_user_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_tags_user_id_index ON public.users_tags USING btree (user_id);


--
-- Name: users_tags_user_id_tag_name_unique_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_tags_user_id_tag_name_unique_index ON public.users_tags USING btree (user_id, tag_name);


--
-- Name: notification notification_send_to_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification
    ADD CONSTRAINT notification_send_to_fkey FOREIGN KEY (send_to) REFERENCES public."user"(id);


--
-- Name: post post_thread_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post
    ADD CONSTRAINT post_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.thread(id);


--
-- Name: post post_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post
    ADD CONSTRAINT post_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: post post_user_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post
    ADD CONSTRAINT post_user_name_fkey FOREIGN KEY (user_name) REFERENCES public."user"(name);


--
-- Name: posts_quotes posts_quotes_quoted_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts_quotes
    ADD CONSTRAINT posts_quotes_quoted_id_fkey FOREIGN KEY (quoted_id) REFERENCES public.post(id);


--
-- Name: posts_quotes posts_quotes_quoter_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts_quotes
    ADD CONSTRAINT posts_quotes_quoter_id_fkey FOREIGN KEY (quoter_id) REFERENCES public.post(id);


--
-- Name: tags_main_tags tags_main_tags_belongs_to_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags_main_tags
    ADD CONSTRAINT tags_main_tags_belongs_to_fkey FOREIGN KEY (belongs_to) REFERENCES public.tag(name);


--
-- Name: tags_main_tags tags_main_tags_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags_main_tags
    ADD CONSTRAINT tags_main_tags_name_fkey FOREIGN KEY (name) REFERENCES public.tag(name);


--
-- Name: thread thread_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thread
    ADD CONSTRAINT thread_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: thread thread_user_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thread
    ADD CONSTRAINT thread_user_name_fkey FOREIGN KEY (user_name) REFERENCES public."user"(name);


--
-- Name: threads_tags threads_tags_tag_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads_tags
    ADD CONSTRAINT threads_tags_tag_name_fkey FOREIGN KEY (tag_name) REFERENCES public.tag(name);


--
-- Name: threads_tags threads_tags_thread_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads_tags
    ADD CONSTRAINT threads_tags_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.thread(id);


--
-- Name: users_tags users_tags_tag_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_tags
    ADD CONSTRAINT users_tags_tag_name_fkey FOREIGN KEY (tag_name) REFERENCES public.tag(name);


--
-- Name: users_tags users_tags_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_tags
    ADD CONSTRAINT users_tags_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- PostgreSQL database dump complete
--

