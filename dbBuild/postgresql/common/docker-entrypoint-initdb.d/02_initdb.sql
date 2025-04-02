-- claimix database is already created in docker-compose and used in current connection;
CREATE TABLESPACE claimix_entities OWNER postgres LOCATION '/var/lib/postgresql/data/claimix_entities';
CREATE TABLESPACE claimix_param OWNER postgres LOCATION '/var/lib/postgresql/data/claimix_param';
CREATE TABLESPACE claimix_oper OWNER postgres LOCATION '/var/lib/postgresql/data/claimix_oper';

ALTER DATABASE claimix SET datestyle TO 'ISO, DMY';