CREATE TABLE property_listing
(
    property_id         BIGINT(20) UNSIGNED        NOT NULL AUTO_INCREMENT,
    user_id             BIGINT(20) UNSIGNED        NOT NULL,
    addresses_id        INT UNSIGNED               NOT NULL,
    listing_description TINYTEXT                   NOT NULL,
    listing_price       INT                        NOT NULL,
    listing_currency    ENUM ('usd', 'hrv', 'eur') NOT NULL,
    listing_is_active   BOOLEAN                    NOT NULL,
    created_at          TIMESTAMP                  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP                  NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (property_id),
    FOREIGN KEY (user_id) REFERENCES seller (user_id),
    FOREIGN KEY (property_id) REFERENCES property (property_id),
    FOREIGN KEY (addresses_id) REFERENCES addresses (addresses_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;