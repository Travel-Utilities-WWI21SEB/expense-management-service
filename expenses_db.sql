-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler version: 1.0.3
-- PostgreSQL version: 15.0
-- Project Site: pgmodeler.io
-- Model Author: ---

-- Database creation must be performed outside a multi lined SQL file. 
-- These commands were put in this file only as a convenience.
-- 
-- object: "travel_expenses-db" | type: DATABASE --
-- DROP DATABASE IF EXISTS "travel_expenses-db";
CREATE SCHEMA "public";
-- ddl-end --


-- object: public."user" | type: TABLE --
-- DROP TABLE IF EXISTS public."user" CASCADE;
CREATE TABLE public."user" (
	id uuid NOT NULL,
	username varchar,
	email varchar,
	password varchar,
	activated boolean,
	CONSTRAINT user_pk PRIMARY KEY (id)
);
-- ddl-end --

-- object: public.trip | type: TABLE --
-- DROP TABLE IF EXISTS public.trip CASCADE;
CREATE TABLE public.trip (
	id uuid NOT NULL,
	location varchar,
	start date,
	"end" date,
	CONSTRAINT travel_pk PRIMARY KEY (id)
);
-- ddl-end --

-- object: public.activation_token | type: TABLE --
-- DROP TABLE IF EXISTS public.activation_token CASCADE;
CREATE TABLE public.activation_token (
	id uuid,
	created_at timestamptz,
	confirmed_at timestamptz,
	id_user uuid NOT NULL,
	CONSTRAINT activation_token_pk PRIMARY KEY (id)
);
-- ddl-end --

-- object: user_fk | type: CONSTRAINT --
-- ALTER TABLE public.activation_token DROP CONSTRAINT IF EXISTS user_fk CASCADE;
ALTER TABLE public.activation_token ADD CONSTRAINT user_fk FOREIGN KEY (id_user)
REFERENCES public."user" (id) MATCH FULL
ON DELETE CASCADE;
-- ddl-end --

-- object: public.many_user_has_many_travel | type: TABLE --
-- DROP TABLE IF EXISTS public.many_user_has_many_travel CASCADE;
CREATE TABLE public.many_user_has_many_travel (
	start date,
	"end" date,
	accepted boolean,
	id_user uuid NOT NULL,
	id_trip uuid NOT NULL,
	CONSTRAINT many_user_has_many_travel_pk PRIMARY KEY (id_user,id_trip)
);
-- ddl-end --

-- object: user_fk | type: CONSTRAINT --
-- ALTER TABLE public.many_user_has_many_travel DROP CONSTRAINT IF EXISTS user_fk CASCADE;
ALTER TABLE public.many_user_has_many_travel ADD CONSTRAINT user_fk FOREIGN KEY (id_user)
REFERENCES public."user" (id) MATCH FULL
ON UPDATE CASCADE;
-- ddl-end --

-- object: trip_fk | type: CONSTRAINT --
-- ALTER TABLE public.many_user_has_many_travel DROP CONSTRAINT IF EXISTS trip_fk CASCADE;
ALTER TABLE public.many_user_has_many_travel ADD CONSTRAINT trip_fk FOREIGN KEY (id_trip)
REFERENCES public.trip (id) MATCH FULL
ON UPDATE CASCADE;
-- ddl-end --

-- object: public.cost_category | type: TABLE --
-- DROP TABLE IF EXISTS public.cost_category CASCADE;
CREATE TABLE public.cost_category (
	id uuid NOT NULL,
	name varchar,
	description varchar,
	icon varchar,
	color smallint,
	id_trip uuid,
	CONSTRAINT cost_category_pk PRIMARY KEY (id)
);
-- ddl-end --

-- object: trip_fk | type: CONSTRAINT --
-- ALTER TABLE public.cost_category DROP CONSTRAINT IF EXISTS trip_fk CASCADE;
ALTER TABLE public.cost_category ADD CONSTRAINT trip_fk FOREIGN KEY (id_trip)
REFERENCES public.trip (id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: public.cost | type: TABLE --
-- DROP TABLE IF EXISTS public.cost CASCADE;
CREATE TABLE public.cost (
	id uuid NOT NULL,
	amount numeric,
	created_at timestamptz,
	deducted_at date,
	"end" date,
	id_cost_category uuid NOT NULL,
	CONSTRAINT cost_pk PRIMARY KEY (id)
);
-- ddl-end --

-- object: cost_category_fk | type: CONSTRAINT --
-- ALTER TABLE public.cost DROP CONSTRAINT IF EXISTS cost_category_fk CASCADE;
ALTER TABLE public.cost ADD CONSTRAINT cost_category_fk FOREIGN KEY (id_cost_category)
REFERENCES public.cost_category (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.many_user_has_many_cost | type: TABLE --
-- DROP TABLE IF EXISTS public.many_user_has_many_cost CASCADE;
CREATE TABLE public.many_user_has_many_cost (
	id_user uuid NOT NULL,
	id_cost uuid NOT NULL,
	is_creditor boolean,
	CONSTRAINT many_user_has_many_cost_pk PRIMARY KEY (id_user,id_cost)
);
-- ddl-end --

-- object: user_fk | type: CONSTRAINT --
-- ALTER TABLE public.many_user_has_many_cost DROP CONSTRAINT IF EXISTS user_fk CASCADE;
ALTER TABLE public.many_user_has_many_cost ADD CONSTRAINT user_fk FOREIGN KEY (id_user)
REFERENCES public."user" (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: cost_fk | type: CONSTRAINT --
-- ALTER TABLE public.many_user_has_many_cost DROP CONSTRAINT IF EXISTS cost_fk CASCADE;
ALTER TABLE public.many_user_has_many_cost ADD CONSTRAINT cost_fk FOREIGN KEY (id_cost)
REFERENCES public.cost (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: public.transaction | type: TABLE --
-- DROP TABLE IF EXISTS public.transaction CASCADE;
CREATE TABLE public.transaction (
	id uuid NOT NULL,
	sender_id uuid,
	receiver_id uuid,
	amount numeric,
	currency varchar,
	CONSTRAINT transaction_pk PRIMARY KEY (id)
);
-- ddl-end --

-- object: debitor_fk | type: CONSTRAINT --
-- ALTER TABLE public.transaction DROP CONSTRAINT IF EXISTS debitor_fk CASCADE;
ALTER TABLE public.transaction ADD CONSTRAINT debitor_fk FOREIGN KEY (sender_id)
REFERENCES public."user" (id) MATCH SIMPLE
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: creditor_fk | type: CONSTRAINT --
-- ALTER TABLE public.transaction DROP CONSTRAINT IF EXISTS creditor_fk CASCADE;
ALTER TABLE public.transaction ADD CONSTRAINT creditor_fk FOREIGN KEY (receiver_id)
REFERENCES public."user" (id) MATCH SIMPLE
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --
