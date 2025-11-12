CREATE TABLE IF NOT EXISTS `school_messenger_users` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `CreatedAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `UpdatedAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `IsActive` tinyint(1) NOT NULL DEFAULT 1,
  `Code` varchar(10) DEFAULT NULL,
  `PSID` varchar(100) NOT NULL,
  `FBName` varchar(100) DEFAULT 'USER',
  `FBImgURL` text DEFAULT NULL,
  `Email` varchar(50) DEFAULT NULL,
  `LastLoginAt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `Notes1` text DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `PSID` (`PSID`),
  UNIQUE KEY `Code` (`Code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_is_active ON school_messenger_users (IsActive);
CREATE INDEX idx_is_registered ON school_messenger_users (IsRegistered);

-- ALTER TABLE gk_miniapps.school_messenger_users ADD IsRegistered TINYINT DEFAULT 0 NOT NULL;
