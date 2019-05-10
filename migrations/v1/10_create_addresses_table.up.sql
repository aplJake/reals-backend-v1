CREATE TABLE addresses
(
    addresses_id  INT UNSIGNED NOT NULL AUTO_INCREMENT,
    city_id       INT UNSIGNED NOT NULL,
    street_name   VARCHAR(150) NOT NULL,
    street_number VARCHAR(10)  NOT NULL,
    PRIMARY KEY (addresses_id),
    FOREIGN KEY (city_id)
        REFERENCES city (city_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;