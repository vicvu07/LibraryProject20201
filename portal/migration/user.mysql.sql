CREATE TABLE `user`
(
  `id` bigint
(20) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar
(32) NOT NULL,
`name` varchar(32),
`role` varchar(10),
`dob` varchar(12),
`sex` varchar(10),
`phonenumber` varchar(15),
 `status` tinyint
(4) NOT NULL DEFAULT '1', 
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `checksum` bigint
(20) unsigned NOT NULL,
  PRIMARY KEY
(`id`),
  UNIQUE KEY `username`
(`username`)
) ENGINE=InnoDB;

CREATE TABLE `user_security`
(
  `username` varchar
(32) NOT NULL,
  `password` varbinary
(60) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `checksum` bigint
(20) unsigned NOT NULL,
  PRIMARY KEY
(`username`)
) ENGINE=InnoDB;
