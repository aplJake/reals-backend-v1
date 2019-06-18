CREATE TABLE notifications
(
    notification_id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    text            TEXT                NOT NULL,
    PRIMARY KEY (notification_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;