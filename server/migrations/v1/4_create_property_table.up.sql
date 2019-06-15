CREATE TABLE property
(
    property_id           BIGINT(20) UNSIGNED                             NOT NULL AUTO_INCREMENT,
    room_number           ENUM ('1 bd', '2 bd', '3 bd', 'more than 3 db') NOT NULL,
    construction_type     SET ('apartment', 'house')                      NOT NULL,
    area                  INT UNSIGNED                                    NOT NULL,
    bathroom_number       TINYINT(20),
    max_floor_number      TINYINT,
    property_floor_number TINYINT,
    PRIMARY KEY (property_id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1;