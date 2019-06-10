CREATE TABLE admins
(
    user_id         BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    admin_role VARCHAR(20)         NOT NULL,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;