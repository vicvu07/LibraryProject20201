drop table if exists borrowform;

create table doc 
(
	`id_borrow` bigint(20) unsigned NOT NULL,
	`id_doc` bigint(20) unsigned NOT NULL,
	`id_cus` bigint(20) unsigned NOT NULL,
	`id_lib` bigint(20) unsigned NOT NULL,
	`status` int,
	`start_at` datetime NOT NULL CURRENT_TIMESTAMP,
	`end_at` datetime NOT NULL,
	`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(`id_borrow`)
) Engine=InnoDB;
