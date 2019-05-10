CREATE TABLE country
(
    country_id   INT UNSIGNED NOT NULL AUTO_INCREMENT,
    country_name VARCHAR(100) NOT NULL,
    zip_code     VARCHAR(8),
    PRIMARY KEY (country_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;