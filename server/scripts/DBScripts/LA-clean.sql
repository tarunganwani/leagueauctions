\c leagueauction

drop index if exists la_schema.IDX_USER_EMAIL		cascade;
drop index if exists la_schema.IDX_PLAYER_ID		cascade;
drop table if exists la_schema.la_user 			cascade;
drop table if exists la_schema.la_player 		cascade;
drop table if exists la_schema.la_user_player_map 	cascade;
drop table if exists la_schema.la_auction_participant 	cascade;
drop table if exists la_schema.la_category		cascade;
drop table if exists la_schema.la_auctionboard		cascade;
drop table if exists la_schema.la_player_pool		cascade;
drop table if exists la_schema.la_participant		cascade;
drop table if exists la_schema.la_purchase		cascade;
drop table if exists la_schema.la_bidderinfo		cascade;
drop table if exists la_schema.la_bid_transaction	cascade;

drop schema if exists la_schema;
\c postgres
drop database if exists leagueauction;


