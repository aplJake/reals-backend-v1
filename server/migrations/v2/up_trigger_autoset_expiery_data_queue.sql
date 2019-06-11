# CREATE TRIGGER test_trigger BEFORE INSERT ON `product`
#     FOR EACH ROW SET
#     NEW.entryDate = IFNULL(NEW.entryDate, NOW()),
#     NEW.expDate = TIMESTAMPADD(DAY, 14, NEW.entryDate);

delimiter #

create trigger auto_set_date_interval_prop_q before insert on reals.property_queue
    for each row

    begin
        set
            NEW.queue_time = DATE_ADD(NOW(), INTERVAL 2 HOUR );
    end #
delimiter ;

insert into property_queue(user_id, property_id) VALUES (1, 1);