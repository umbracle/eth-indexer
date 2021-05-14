/*
CREATE TABLE IF NOT EXISTS events (
    event       text,
	indx 		numeric,
	full_indx   numeric,
	tx_index 	numeric,
	tx_hash 	text,
	block_num 	numeric,
	block_hash 	text,
	address 	text,
	topicid     text,
	topics 		text,
	data 		text,
    removed     boolean
);

CREATE INDEX IF NOT EXISTS events_block_num ON events(block_num);
*/