BEGIN;

DROP TABLE IF EXISTS public.users ;
DROP TABLE IF EXISTS public.actors cascade ;
DROP TABLE IF EXISTS public.films cascade ;
DROP TABLE IF EXISTS public.films_actors ;

END;