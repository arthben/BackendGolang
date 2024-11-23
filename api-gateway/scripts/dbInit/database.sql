CREATE EXTENSION postgis;

create table rideindego_master (
	fetch_id uuid PRIMARY KEY,
	type_collection varchar(25) not NULL,
	last_update TIMESTAMP WITH TIME zone not NULL
);


create table rideindego_features (
	fetch_id uuid not null,
	feat_id integer not null,
	ftype varchar(25) not null,
	geo_type varchar(10) not null,
	geo_coordinate GEOMETRY(Point,4326),
	primary key(fetch_id, feat_id)
);


create table rideindego_properties(
	fetch_id uuid not null,
	feat_id integer not null,
	id integer not null,
	name varchar,
	notes varchar default null,
	kiosk_id integer,
	event_end varchar default null,
	latitude float,
	open_time varchar default null,
	time_zone varchar default null,
	close_time varchar default null,
	is_virtual bool,
	kiosk_type integer,
	longitude float,
	event_start varchar default null,
	public_text varchar,
	total_docks integer,
	address_city varchar,
	coordinates GEOMETRY(Point,4326),
	kiosk_status varchar,
	address_state varchar,
	is_event_based bool,
	address_street varchar,
	address_zip_code varchar,
	bikes_available integer,
	docks_available integer,
	trikes_available integer,
	kiosk_public_status varchar,
	smart_bikes_available integer,
	reward_bikes_available integer,
	reward_docks_available integer,
	classic_bikes_available integer,
	kiosk_connection_status varchar,
	electric_bikes_available integer,
	primary key(fetch_id, feat_id, id)
);
create index idx_kioskid on rideindego_properties(kiosk_id);

create table rideindego_properties_bikes(
	fetch_id uuid not null,
	feat_id integer not null,
	id integer not null,
	idx SERIAL not null,
	battery INTEGER,
	dock_number INTEGER,
	is_electric bool,
	is_available bool,
	primary key(fetch_id, feat_id, id, idx)
);

create table openweather_master (
	fetch_id uuid not null,
	base varchar,
	clouds integer,
	cod integer,
	coord GEOMETRY(Point,4326),
	dt integer,
	id integer,
	main_feels_like float,
	main_grnd_level integer,
	main_humidity integer,
	main_pressure integer,
	main_sea_level integer,
	main_temp float,
	main_temp_max float,
	main_temp_min float,
	name varchar,
	rain_one_hour float,
	sys_country varchar,
	sys_id integer,
	sys_sunrise integer,
	sys_sunset integer,
	sys_type integer,
	timezone integer,
	visibility integer,
	wind_deg integer,
	wind_speed float,
	primary key(fetch_id)
);

create table openweather_weather(
	fetch_id uuid not null,
	idx serial not null,
	description varchar,
	icon varchar,
	id integer,
	main varchar,
	primary key(fetch_id, idx)
);