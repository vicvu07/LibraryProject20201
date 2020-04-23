CREATE TABLE `user`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(32) NOT NULL,
  `pwd_counter` int(11) NOT NULL DEFAULT 5, 
  `status` tinyint(4) NOT NULL DEFAULT '1', 
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `checksum` bigint(20) unsigned NOT NULL,
  PRIMARY KEY(`id`),
  UNIQUE KEY `username`(`username`)
) ENGINE=InnoDB;
 
CREATE TABLE `user_security`
(
  `username` varchar(32) NOT NULL,
  `password` varbinary(60) NOT NULL, 
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `checksum` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`username`)
) ENGINE=InnoDB;

CREATE TABLE `user_detail`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL,
  `DOB` varchar(10),
  `sex` varchar(7),
  `position` varchar(40),
  `phonenum` varchar(15) NOT NULL,
  `national_id` varchar(20) NOT NULL,
  `salary` bigint(20),
  `username` varchar(32) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE (`username`)
) ENGINE=InnoDB;

CREATE TABLE `department`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL,
  `description` varchar(200),
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,  
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `total_salary` bigint(20),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

CREATE TABLE `department_management` 
(
  `department_id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `status` int(6) NOT NULL DEFAULT 1 COMMENT '1 is active, 0 is not active', 
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,  
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (`department_id`,`user_id`)
) ENGINE=InnoDB;

CREATE TABLE `plan` 
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL,
  `description` varchar(200) NOT NULL,
  `father_plan_id` bigint(20) DEFAULT 0,
  `current_status` varchar(1) NOT NULL COMMENT 'D is done, O is not on going', 
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,  
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

CREATE TABLE `plan_management`
(
  `IndiOrDepart` varchar(1) NOT NULL,
  `foreign_id` bigint(20) NOT NULL,
  `plan_id` bigint(20) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,  
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (`IndiOrDepart`,`foreign_id`,`plan_id`)
) ENGINE=InnoDB;