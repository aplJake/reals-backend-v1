CREATE TABLE users
(
    user_id       BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    user_name     VARCHAR(100)        NOT NULL,
    user_email    VARCHAR(50)         NOT NULL,
    user_password VARCHAR(255)        NOT NULL,
    PRIMARY KEY (user_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;