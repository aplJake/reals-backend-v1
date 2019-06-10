CREATE TABLE city
(
    city_id    INT UNSIGNED NOT NULL AUTO_INCREMENT,
    country_id INT UNSIGNED NOT NULL,
    city_name  VARCHAR(100) NOT NULL,
    PRIMARY KEY (city_id),
    FOREIGN KEY (country_id)
        REFERENCES country (country_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;