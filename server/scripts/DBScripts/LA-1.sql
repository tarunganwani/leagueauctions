create database leagueauction;

\c leagueauction

create schema la_schema;

/*postgres extension to use uuid datatype*/
create extension IF NOT EXISTS "uuid-ossp";


create table IF NOT EXISTS la_schema.la_user (
	user_id 		UUID 	DEFAULT uuid_generate_v4 () NOT NULL,
	email_id		TEXT 	UNIQUE	NOT NULL,
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
	player_type		INT,
	player_photo		BYTEA,
	is_active	BOOLEAN 	DEFAULT TRUE NOT NULL,
	PRIMARY KEY(player_id)
);

/* create index on primary key to facilicate use of UPSERT (INSERT.. ON CONFLICT) query*/
create index IDX_PLAYER_ID on
la_schema.la_player(player_id);

create table IF NOT EXISTS la_schema.la_user_player_map (
	map_id		UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	user_id		UUID	UNIQUE,
	player_id	UUID	UNIQUE,
	PRIMARY KEY(map_id),
	FOREIGN KEY(user_id)	REFERENCES la_schema.la_user(user_id) ON DELETE CASCADE,
	FOREIGN KEY(player_id)	REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE
);


/*Auction participant
create table IF NOT EXISTS la_schema.la_auction_participant (
	participant_id		UUID 		DEFAULT uuid_generate_v4 () 	NOT NULL,
	participant_role	INT			NOT NULL,
	player_id			UUID		NOT NULL,
	PRIMARY KEY(participant_id),
	FOREIGN KEY(player_id) REFERENCES la_schema.la_player(player_id)
);*/


/*
	Auction board table
	Generate a random 6 digit auction code
*/
create table IF NOT EXISTS la_schema.la_auctionboard (
	auction_id		UUID 		DEFAULT uuid_generate_v4 ()  		NOT NULL,
	auctioneer_id	UUID		NOT NULL,
	auction_name	TEXT		NOT NULL,
	schedule_time	TIMESTAMPTZ	NOT NULL,
	purse			BIGINT		NOT NULL,
	purse_ccy		TEXT		DEFAULT 'COINS' NOT NULL,
	is_active		BOOLEAN 	DEFAULT TRUE NOT NULL,
	auction_code	INT			DEFAULT floor(random()* (999999-100000 + 1) + 100000)	NOT NULL,
	PRIMARY KEY(auction_id),
	FOREIGN KEY(auctioneer_id) REFERENCES la_schema.la_player(player_id)
);

/*Player category*/
create table IF NOT EXISTS la_schema.la_category (
	category_id		UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	auction_id		UUID	NOT NULL,
	category_name	TEXT	NOT NULL,
	base_price		BIGINT	NOT NULL,
	PRIMARY KEY(category_id),
	FOREIGN KEY(auction_id)		REFERENCES la_schema.la_auctionboard(auction_id) ON DELETE CASCADE

);

/*Players going under the hammer*/
/*create table IF NOT EXISTS la_schema.la_player_pool(
	pool_id			UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	auction_id		UUID	NOT NULL,
	player_id		UUID	NOT NULL,
	category_id		UUID	NOT NULL,
	plater_status	INT	NOT NULL,
	PRIMARY KEY(pool_id),
	FOREIGN KEY(auction_id)		REFERENCES la_schema.la_auctionboard(auction_id) ON DELETE CASCADE,
	FOREIGN KEY(category_id)	REFERENCES la_schema.la_category(category_id) ON DELETE CASCADE,
	FOREIGN KEY(player_id)		REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE
);*/

/*Participants involved in an auction*/
create table IF NOT EXISTS la_schema.la_participant(
	participant_id	UUID 	DEFAULT uuid_generate_v4 () 	NOT NULL,
	player_id	UUID	NOT NULL,
	auction_id	UUID	NOT NULL,
	participant_role	INT	NOT NULL,
	category_id			UUID	NOT NULL,
	PRIMARY KEY(participant_id),
	FOREIGN KEY(auction_id)	REFERENCES la_schema.la_auctionboard(auction_id) ON DELETE CASCADE,
	FOREIGN KEY(player_id)	REFERENCES la_schema.la_player(player_id) ON DELETE CASCADE
);

comment on column la_schema.la_participant.participant_id is 'Participant UUID';
comment on column la_schema.la_participant.player_id is 'Associated player UUID';
comment on column la_schema.la_participant.auction_id is 'Associated auction board UUID';
comment on column la_schema.la_participant.participant_role is 'Role: 1-Captain(Bidder) 2-Player 3-Viewer ';
comment on column la_schema.la_participant.category_id is 'Category UUID';

-- Bidding tables

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


CREATE TYPE category_item AS (
    category_id    			UUID,
    auction_id     			UUID,
    category_name           TEXT,
    base_price           		BIGINT
);

-- Functions

/*
create function dummy_fn()
returns int
language plpgsql
as
$$
declare
	out_auction_code integer;
begin
	return 1;
end;
$$;*/

create function insert_auction_board_info(in_auction_id UUID, 
										in_auctioneer_id UUID, 
										in_auction_name TEXT, 
										in_schedule_time TIMESTAMPTZ, 
										in_purse BIGINT,
										in_purse_ccy TEXT,
										in_category_item_list category_item[])
returns int
language plpgsql
as
$$
declare
	out_auction_code INT;
	in_category_item category_item;
begin
	INSERT INTO la_schema.la_auctionboard(auction_id, auctioneer_id, auction_name, schedule_time, purse, purse_ccy)
	VALUES(in_auction_id, in_auctioneer_id, in_auction_name, in_schedule_time, in_purse, in_purse_ccy)
	RETURNING auction_code into out_auction_code;
	FOREACH in_category_item IN ARRAY in_category_item_list LOOP
		INSERT INTO la_schema.la_category(category_id, auction_id, category_name, base_price)
		VALUES(in_category_item.category_id, in_category_item.auction_id, in_category_item.category_name, in_category_item.base_price);
	END LOOP;
	return out_auction_code;
end;
$$;