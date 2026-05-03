
create table miinstance (
	miinstance_passport varchar(10) primary key,
	miinstance_name varchar(1000),
	miinstance_type varchar(1000),
	miinstance_state_condition varchar(40),
	miinstance_tech_condition varchar(40),
	issue_date date null,
	commissioning_date date null,
	is_fit bool,
	MPI int
)with(fillfactor=85);


CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_miinstance_name_trgm 
ON miinstance 
USING GIN(miinstance_name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_miinstance_passport_trgm 
ON miinstance 
USING GIN(miinstance_passport gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_miinstance_type_trgm 
ON miinstance 
USING GIN(miinstance_type gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_miinstance_state_condition 
ON miinstance(miinstance_state_condition) 
WITH (fillfactor = 85);

CREATE INDEX IF NOT EXISTS idx_miinstance_tech_condition 
ON miinstance(miinstance_tech_condition) 
WITH (fillfactor = 85);

CREATE INDEX IF NOT EXISTS idx_miinstance_issue_date 
ON miinstance(issue_date) 
WITH (fillfactor = 85);

CREATE INDEX IF NOT EXISTS idx_miinstance_commissioning_date 
ON miinstance(commissioning_date) 
WITH (fillfactor = 85);

CREATE INDEX IF NOT EXISTS idx_miinstance_mpi 
ON miinstance(MPI) 
WHERE MPI IS NOT NULL;