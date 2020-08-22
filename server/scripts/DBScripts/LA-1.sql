create database leagueauction;

\c leagueauction

create schema la_schema;

create extension IF NOT EXISTS "uuid-ossp";

create table IF NOT EXISTS la_schema.la_user (
	user_id 	UUID 	DEFAULT uuid_generate_v4 () NOT NULL,
	email_id	TEXT 	UNIQUE	NOT NULL,
	password_hash	TEXT 	NOT NULL,
	password_salt	TEXT 	NOT NULL,
	activation_code	TEXT,
	is_active	BOOLEAN 	DEFAULT FALSE NOT NULL,

	PRIMARY KEY(user_id)
);	

create index IDX_USER_EMAIL on
la_schema.la_user(email_id);

create table IF NOT EXISTS la_schema.la_player (
	player_id		UUID 	DEFAULT uuid_generate_v4 () NOT NULL,
	player_name		TEXT	NOT NULL,
	player_bio		TEXT,
	player_profile_link	TEXT,
	player_type		TEXT,
	player_photo		BYTEA,
	is_active	BOOLEAN 	DEFAULT TRUE NOT NULL,
	PRIMARY KEY(player_id)
);

create table IF NOT EXISTS la_schema.la_user_player_map (
	map_id		UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	user_id		UUID,
	player_id	UUID,
	PRIMARY KEY(map_id),
	FOREIGN KEY(user_id)	REFERENCES la_schema.la_user(user_id) ON DELETE CASCADE,
	FOREIGN KEY(player_id)	REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE
);

create table IF NOT EXISTS la_schema.la_auction_participant (
	participant_id	UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	role		TEXT	NOT NULL,
	player_id	UUID	NOT NULL,
	PRIMARY KEY(participant_id),
	FOREIGN KEY(player_id) REFERENCES la_schema.la_player(player_id)
);


create table IF NOT EXISTS la_schema.la_category (
	category_id	UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	category_name	TEXT	NOT NULL,
	base_price	MONEY	NOT NULL,
	PRIMARY KEY(category_id)
);

create table IF NOT EXISTS la_schema.la_auctionboard (
	auction_id	UUID 	DEFAULT uuid_generate_v4 ()  		NOT NULL,
	auctioneer_id	UUID		NOT NULL,
	schedule_time	TIMESTAMPTZ	NOT NULL,
	is_active	BOOLEAN 	DEFAULT FALSE NOT NULL,
	PRIMARY KEY(auction_id),
	FOREIGN KEY(auctioneer_id) REFERENCES la_schema.la_player(player_id)
);

create table IF NOT EXISTS la_schema.la_player_pool(
	pool_id		UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	auction_id	UUID	NOT NULL,
	player_id	UUID	NOT NULL,
	category_id	UUID	NOT NULL,
	status		TEXT	NOT NULL,
	PRIMARY KEY(pool_id),
	FOREIGN KEY(auction_id)		REFERENCES la_schema.la_auctionboard(auction_id) ON DELETE CASCADE,
	FOREIGN KEY(category_id)	REFERENCES la_schema.la_category(category_id) ON DELETE CASCADE,
	FOREIGN KEY(player_id)		REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE
);


create table IF NOT EXISTS la_schema.la_participant(
	participant_id	UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	player_id	UUID	NOT NULL,
	auction_id	UUID	NOT NULL,
	role		TEXT	NOT NULL,
	PRIMARY KEY(participant_id),
	FOREIGN KEY(auction_id)	REFERENCES la_schema.la_auctionboard(auction_id) ON DELETE CASCADE,
	FOREIGN KEY(player_id)	REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE
);

create table IF NOT EXISTS la_schema.la_purchase (
	purchase_id	UUID 	DEFAULT uuid_generate_v4 ()  	NOT NULL,
	player_id	UUID 	NOT NULL,
	amount		MONEY	NOT NULL,
	PRIMARY KEY(purchase_id),
	FOREIGN KEY(player_id)	REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE
);

create table IF NOT EXISTS la_schema.la_bidderinfo (
	bidderinfo_id 	UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	auction_id	UUID 	NOT NULL,
	bidder_id	UUID 	NOT NULL,
	current_purse	MONEY	NOT NULL,
	purchase_id	UUID	NOT NULL,
	PRIMARY KEY(bidder_id),
	FOREIGN KEY(auction_id) REFERENCES la_schema.la_auctionboard(auction_id) ON DELETE CASCADE,
	FOREIGN KEY(purchase_id) REFERENCES la_schema.la_purchase(purchase_id) ON DELETE CASCADE,
	FOREIGN KEY(bidder_id)	REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE
);

create table IF NOT EXISTS la_schema.la_bid_transaction (
	bid_id		UUID 	DEFAULT uuid_generate_v4 ()  		NOT NULL,
	auction_id	UUID		NOT NULL,
	bidder_id	UUID		NOT NULL,
	player_id	UUID 		NOT NULL,
	bid_amount	MONEY		NOT NULL,
	bid_timestamp	TIMESTAMP	NOT NULL,
	PRIMARY KEY(bid_id),
	FOREIGN KEY(auction_id) REFERENCES la_schema.la_auctionboard(auction_id) ON DELETE CASCADE,
	FOREIGN KEY(bidder_id)	REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE,
	FOREIGN KEY(player_id)	REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE
);
