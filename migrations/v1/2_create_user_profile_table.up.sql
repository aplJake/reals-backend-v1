CREATE TABLE user_profile
(
    user_id             BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    profile_description TINYTEXT,
    created_at          TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id),
    CONSTRAINT fk_user_profile
        FOREIGN KEY (user_id)
            REFERENCES users (user_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;