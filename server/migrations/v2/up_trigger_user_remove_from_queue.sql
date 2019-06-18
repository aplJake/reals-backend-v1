delimiter #

create trigger notification_on_remove_from_q
    after delete
    on reals.property_queue
    for each row

begin
    declare property_owner_id bigint;
    select user_id from property_listing where property_listing.property_id = OLD.property_id into property_owner_id;
    insert into notifications(user_id, text)
    values (property_owner_id, 'User was removed from your property listing queue');
end#
delimiter ;