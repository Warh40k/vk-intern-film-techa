BEGIN;

ALTER table films_actors DROP CONSTRAINT films_actors_pkey;
ALTER table films_actors ADD CONSTRAINT films_actors_pkey PRIMARY KEY (film_id, actor_id);
ALTER table films_actors DROP COLUMN id;

END;