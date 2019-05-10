CREATE TABLE seller
(
    user_id          BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    telephone_number VARCHAR(12),
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES user_profile (user_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;