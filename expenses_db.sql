-- object: public | type: SCHEMA --
DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;
-- ddl-end --

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- ddl-end --

-- object: public."user" | type: TABLE --
DROP TABLE IF EXISTS public."user" CASCADE;
CREATE TABLE public."user"
(
    id        uuid              NOT NULL DEFAULT uuid_generate_v4(),
    username  character varying NOT NULL,
    email     character varying NOT NULL,
    password  character varying NOT NULL,
    activated boolean           NOT NULL,
    CONSTRAINT user_pk PRIMARY KEY (id),
    CONSTRAINT username_un UNIQUE (username),
    CONSTRAINT email_un UNIQUE (email)
);
-- ddl-end --

-- object: public.trip | type: TABLE --
DROP TABLE IF EXISTS public.trip CASCADE;
CREATE TABLE public.trip
(
    id         uuid NOT NULL DEFAULT uuid_generate_v4(),
    name       character varying,
    description character varying,
    location   character varying,
    start_date date NOT NULL,
    end_date   date NOT NULL,
    CONSTRAINT travel_pk PRIMARY KEY (id)
);
-- ddl-end --

-- object: public.token | type: TABLE --
DROP TABLE IF EXISTS public.token CASCADE;
CREATE TABLE public.token
(
    id           uuid NOT NULL DEFAULT uuid_generate_v4(),
    id_user      uuid NOT NULL,
    token        character varying,
    type         character varying,
    created_at   timestamp with time zone,
    confirmed_at timestamp with time zone,
    expires_at   timestamp with time zone,
    CONSTRAINT token_pk PRIMARY KEY (id),
    CONSTRAINT token_un UNIQUE (token)
);
-- ddl-end --

-- object: public.user_trip_association | type: TABLE --
DROP TABLE IF EXISTS public.user_trip_association CASCADE;
CREATE TABLE public.user_trip_association
(
    presence_start_date date    NOT NULL,
    presence_end_date   date    NOT NULL,
    is_accepted         boolean NOT NULL,
    id_user             uuid    NOT NULL,
    id_trip             uuid    NOT NULL,
    CONSTRAINT many_user_has_many_travel_pk PRIMARY KEY (id_user, id_trip)
);
-- ddl-end --

-- object: public.cost_category | type: TABLE --
DROP TABLE IF EXISTS public.cost_category CASCADE;
CREATE TABLE public.cost_category
(
    id          uuid NOT NULL DEFAULT uuid_generate_v4(),
    name        character varying,
    description character varying,
    icon        character varying,
    color       character varying,
    id_trip     uuid NOT NULL,
    CONSTRAINT cost_category_pk PRIMARY KEY (id),
    CONSTRAINT cost_category_un UNIQUE (name, id_trip)
);
-- ddl-end --

-- object: public.cost | type: TABLE --
DROP TABLE IF EXISTS public.cost CASCADE;
CREATE TABLE public.cost
(
    id               uuid NOT NULL DEFAULT uuid_generate_v4(),
    amount           numeric,
    description      character varying,
    created_at       timestamp with time zone,
    deducted_at      date,
    end_date         date,
    id_cost_category uuid NOT NULL,
    CONSTRAINT cost_pk PRIMARY KEY (id)
);
-- ddl-end --

-- object: public.user_cost_association | type: TABLE --
DROP TABLE IF EXISTS public.user_cost_association CASCADE;
CREATE TABLE public.user_cost_association
(
    id_user     uuid    NOT NULL,
    id_cost     uuid    NOT NULL,
    is_creditor boolean NOT NULL,
    amount      numeric NOT NULL,
    CONSTRAINT many_user_has_many_cost_pk PRIMARY KEY (id_user, id_cost)
);
-- ddl-end --

-- object: public.transaction | type: TABLE --
DROP TABLE IF EXISTS public.transaction CASCADE;
CREATE TABLE public.transaction
(
    id          uuid NOT NULL DEFAULT uuid_generate_v4(),
    sender_id   uuid,
    receiver_id uuid,
    amount      numeric,
    currency    character varying,
    CONSTRAINT transaction_pk PRIMARY KEY (id)
);
-- ddl-end --

-- object: user_fk | type: CONSTRAINT --
-- ALTER TABLE public.token DROP CONSTRAINT IF EXISTS user_fk CASCADE;
ALTER TABLE public.token
    ADD CONSTRAINT user_fk FOREIGN KEY (id_user)
        REFERENCES public."user" (id) MATCH FULL
        ON DELETE CASCADE ON UPDATE NO ACTION;
-- ddl-end --

-- object: user_fk | type: CONSTRAINT --
-- ALTER TABLE public.user_trip_association DROP CONSTRAINT IF EXISTS user_fk CASCADE;
ALTER TABLE public.user_trip_association
    ADD CONSTRAINT user_fk FOREIGN KEY (id_user)
        REFERENCES public."user" (id) MATCH FULL
        ON DELETE NO ACTION ON UPDATE CASCADE;
-- ddl-end --

-- object: trip_fk | type: CONSTRAINT --
-- ALTER TABLE public.user_trip_association DROP CONSTRAINT IF EXISTS trip_fk CASCADE;
ALTER TABLE public.user_trip_association
    ADD CONSTRAINT trip_fk FOREIGN KEY (id_trip)
        REFERENCES public.trip (id) MATCH FULL
        ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: trip_fk | type: CONSTRAINT --
-- ALTER TABLE public.cost_category DROP CONSTRAINT IF EXISTS trip_fk CASCADE;
ALTER TABLE public.cost_category
    ADD CONSTRAINT trip_fk FOREIGN KEY (id_trip)
        REFERENCES public.trip (id) MATCH FULL
        ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: cost_category_fk | type: CONSTRAINT --
-- ALTER TABLE public.cost DROP CONSTRAINT IF EXISTS cost_category_fk CASCADE;
ALTER TABLE public.cost
    ADD CONSTRAINT cost_category_fk FOREIGN KEY (id_cost_category)
        REFERENCES public.cost_category (id) MATCH FULL
        ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: user_fk | type: CONSTRAINT --
-- ALTER TABLE public.user_cost_association DROP CONSTRAINT IF EXISTS user_fk CASCADE;
ALTER TABLE public.user_cost_association
    ADD CONSTRAINT user_fk FOREIGN KEY (id_user)
        REFERENCES public."user" (id) MATCH FULL
        ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: cost_fk | type: CONSTRAINT --
-- ALTER TABLE public.user_cost_association DROP CONSTRAINT IF EXISTS cost_fk CASCADE;
ALTER TABLE public.user_cost_association
    ADD CONSTRAINT cost_fk FOREIGN KEY (id_cost)
        REFERENCES public.cost (id) MATCH FULL
        ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: debitor_fk | type: CONSTRAINT --
-- ALTER TABLE public.transaction DROP CONSTRAINT IF EXISTS debitor_fk CASCADE;
ALTER TABLE public.transaction
    ADD CONSTRAINT debitor_fk FOREIGN KEY (sender_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: creditor_fk | type: CONSTRAINT --
-- ALTER TABLE public.transaction DROP CONSTRAINT IF EXISTS creditor_fk CASCADE;
ALTER TABLE public.transaction
    ADD CONSTRAINT creditor_fk FOREIGN KEY (receiver_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: "grant_CU_26541e8cda" | type: PERMISSION --
GRANT CREATE, USAGE
    ON SCHEMA public
    TO pg_database_owner;
-- ddl-end --

-- object: "grant_U_cd8e46e7b6" | type: PERMISSION --
GRANT USAGE
    ON SCHEMA public
    TO PUBLIC;
-- ddl-end --


