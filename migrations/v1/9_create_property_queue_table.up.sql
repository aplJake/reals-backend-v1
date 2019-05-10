CREATE TABLE property_queue
(
    user_id     BIGINT(20) UNSIGNED NOT NULL,
    property_id BIGINT(20) UNSIGNED NOT NULL,
    FOREIGN KEY (user_id) REFERENCES buyer (user_id),
    FOREIGN KEY (property_id) REFERENCES property_listing (property_id)
) ENGINE = InnoDB;