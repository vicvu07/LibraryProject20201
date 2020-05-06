drop table if exists doc;

create table doc 
(
	`id_doc` bigint(20) unsigned NOT NULL,
	`doc_name` varchar(100),
	`doc_author` varchar(30),
	`doc_type` varchar(30),
	`doc_description` varchar(100),
	`status` int,
	`id_borrow` bigint(20) unsigned,
	`fee` bigint(20),
	`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(`id_doc`)
) Engine=InnoDB;
