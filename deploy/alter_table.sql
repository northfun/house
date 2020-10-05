 CREATE TABLE table_cosmetics (     
	id SERIAL,      
	brand character varying(64),      
	name character varying(128),
  	e_name character varying(256),
  	volume character varying(128),
	price character varying(128),
  	type smallint,
  	grassed smallint,
  	purchased smallint,
  	like_rate smallint,
  	effect smallint[],
  	goods_id integer,
	ingredient_ids smallint[]
);

CREATE TABLE table_users (     
	id SERIAL,
  	name character varying(64),
  	comment_num smallint,
  	grassed_num smallint,
  	purchased_num smallint,
  	annual_consume_l integer,
  	annual_consume_r integer,
  	level smallint,
  	uid integer
);

CREATE TABLE table_comments (
   	cosmetic_id integer,
  	star smallint,
  	skin_type smallint,
  	age smallint,
  	comment_length smallint,
  	date character varying(64),
  	uid integer,
  	comment text
);

CREATE TABLE table_ingredients(
   	id SERIAL,
  	extra_names varchar(64)[],
	main_name varchar(64),
	brief varchar(64)[],
  	info text
);

CREATE TABLE table_effect (
   	id SERIAL,
  	effect_name character varying(64)
);


CREATE INDEX index_cosmetics ON table_cosmetics(id,goods_id);
CREATE INDEX index_users ON table_users(id,uid);
CREATE INDEX index_comments ON table_comments(cosmetic_id,uid);
CREATE INDEX index_ingredients ON table_ingredients(id,main_name);

insert into table_users_bk (select * from table_users a where a.id in (select min(id) from table_users group by name,comment_num,grassed_num,purchased_num,annual_consume_l,annual_consume_r, level,uid));

select count(*) from table_users a where a.id in (select min(id) from table_users group by name,comment_num,grassed_num,purchased_num,annual_consume_l,annual_consume_r, level,uid);

select * from (select count(*) as cnum,uid from table_users_bk group by uid) as a where a.cnum > 1;

alter table table_cosmetics add column find_name varchar(256);
