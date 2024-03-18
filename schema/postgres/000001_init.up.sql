BEGIN;

CREATE TABLE IF NOT EXISTS public.users
(
    id serial primary key,
    username character varying(255) NOT NULL UNIQUE ,
    password_hash character varying NOT NULL,
    role smallint NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS public.actors
(
    id serial primary key,
    name character varying(255) NOT NULL,
    birthday date NOT NULL,
    gender smallint NOT NULL
);

CREATE TABLE IF NOT EXISTS public.films
(
    id serial primary key,
    title character varying(150) NOT NULL,
    description character varying(1000) NOT NULL,
    released date NOT NULL,
    rating smallint NOT NULL
);

CREATE TABLE IF NOT EXISTS public.films_actors
(
    id serial primary key ,
    actor_id int references actors(id) on delete cascade,
    film_id int references films(id) on delete cascade
);

END;