--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: posts; Type: TABLE; Schema: public; Owner: rasil
--

CREATE TABLE public.posts (
    id character varying(36) NOT NULL,
    title character varying(255) NOT NULL,
    article_content jsonb,
    small_description character varying(255),
    image text,
    slug character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    user_id character varying(35),
    site_id character varying(36)
);


ALTER TABLE public.posts OWNER TO rasil;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: rasil
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO rasil;

--
-- Name: sites; Type: TABLE; Schema: public; Owner: rasil
--

CREATE TABLE public.sites (
    id character varying(36) NOT NULL,
    name character varying(35) NOT NULL,
    description character varying(150),
    subdirectory character varying(40) NOT NULL,
    image_url text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    user_id character varying(35)
);


ALTER TABLE public.sites OWNER TO rasil;

--
-- Name: subscriptions; Type: TABLE; Schema: public; Owner: rasil
--

CREATE TABLE public.subscriptions (
    stripe_subscription_id text NOT NULL,
    "interval" character varying(100) NOT NULL,
    status character varying(100) NOT NULL,
    plan_id text NOT NULL,
    current_period_start integer NOT NULL,
    current_period_end integer NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    user_id character varying(35)
);


ALTER TABLE public.subscriptions OWNER TO rasil;

--
-- Name: users; Type: TABLE; Schema: public; Owner: rasil
--

CREATE TABLE public.users (
    id character varying(35) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    profile_image text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    customer_id text
);


ALTER TABLE public.users OWNER TO rasil;

--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: posts posts_user_id_site_id_slug_key; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_user_id_site_id_slug_key UNIQUE (user_id, site_id, slug);


--
-- Name: schema_migration schema_migration_pkey; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.schema_migration
    ADD CONSTRAINT schema_migration_pkey PRIMARY KEY (version);


--
-- Name: sites sites_pkey; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.sites
    ADD CONSTRAINT sites_pkey PRIMARY KEY (id);


--
-- Name: sites sites_subdirectory_key; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.sites
    ADD CONSTRAINT sites_subdirectory_key UNIQUE (subdirectory);


--
-- Name: subscriptions subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_pkey PRIMARY KEY (stripe_subscription_id);


--
-- Name: subscriptions subscriptions_user_id_key; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_user_id_key UNIQUE (user_id);


--
-- Name: users users_customer_id_key; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_customer_id_key UNIQUE (customer_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: rasil
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: posts posts_sites_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_sites_id_fk FOREIGN KEY (site_id) REFERENCES public.sites(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: posts posts_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: sites sites_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.sites
    ADD CONSTRAINT sites_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: subscriptions subscriptions_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: rasil
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

