CREATE EVENT delete_event
    ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 10 MINUTE 
    COMMENT '2H interval autoremove from queue'

    DO BEGIN
    DELETE FROM property_queue WHERE queue_time < DATE_SUB(NOW(), INTERVAL 2 HOUR );
END;