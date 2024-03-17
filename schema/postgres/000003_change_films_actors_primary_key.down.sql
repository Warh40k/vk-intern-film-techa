BEGIN;

ALTER table films_actors DROP CONSTRAINT films_actors_pkey;
ALTER table films_actors ADD COLUMN id serial PRIMARY KEY ;

END;