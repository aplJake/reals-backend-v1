# CREATE TABLE profile_notification
# (
#     notification_id BIGINT(20) UNSIGNED NOT NULL,
#     user_id         BIGINT(20) UNSIGNED NOT NULL,
#     FOREIGN KEY (user_id) REFERENCES users (user_id),
#     FOREIGN KEY (notification_id) REFERENCES notification (notification_id)
# ) ENGINE = InnoDB;

CREATE TABLE notifications
(
    notification_id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id         BIGINT(20) UNSIGNED NOT NULL,
    text            TEXT                NOT NULL,
    PRIMARY KEY (notification_id),
    FOREIGN KEY (user_id) REFERENCES user_profile (user_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;
