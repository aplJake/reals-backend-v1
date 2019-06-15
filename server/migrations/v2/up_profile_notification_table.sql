CREATE TABLE profile_notification
(
    notification_id     BIGINT(20) UNSIGNED NOT NULL,
    user_id BIGINT(20) UNSIGNED NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id),
    FOREIGN KEY (notification_id) REFERENCES notification (notification_id)
) ENGINE = InnoDB;

CREATE TABLE notification
(
    notification_id      BIGINT(20) UNSIGNED       NOT NULL AUTO_INCREMENT,
    text TEXT NOT NULL,
    PRIMARY KEY (notification_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;
