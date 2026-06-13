
-- migrate: up
CREATE TABLE IF NOT EXISTS public.plan
(
    plan_uuid uuid NOT NULL DEFAULT uuidv7(),
    plan_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    plan_start_date timestamp without time zone NOT NULL,
    plan_end_date timestamp without time zone NOT NULL,
    CONSTRAINT plan_pkey PRIMARY KEY (plan_uuid)
)

WITH (
    FILLFACTOR = 85
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.plan
    OWNER to postgres;

CREATE INDEX IF NOT EXISTS plan_end_date_index
    ON public.plan USING btree
    (plan_end_date ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS plan_start_date_index
    ON public.plan USING btree
    (plan_start_date ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS public.miinstance_for_plan (
    miinstance_uuid uuid NOT NULL,
    plan_uuid uuid NOT NULL,
    miinstance_name character varying(500) NOT NULL,
    miinstance_type character varying(100) NOT NULL,
    type_of_action actions NOT NULL DEFAULT 'consider'::actions,
    CONSTRAINT miinstance_for_plan_plan_uuid_fkey FOREIGN KEY (plan_uuid)
        REFERENCES public.plan (plan_uuid) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS miinstance_name_for_plan
    ON public.miinstance_for_plan USING btree (plan_uuid ASC NULLS LAST);

CREATE INDEX IF NOT EXISTS plan_uuid_for_plan
    ON public.miinstance_for_plan USING btree (plan_uuid ASC NULLS LAST);

-- migrate: down

DROP TABLE IF EXISTS public.miinstance_for_plan;
DROP TABLE IF EXISTS public.plan;