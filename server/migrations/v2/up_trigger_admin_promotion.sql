delimiter #

create trigger admin_promotion
    after insert
    on reals.admins
    for each row
begin
    insert into notifications(user_id, text) values (new.user_id, 'You have been promoted to manager user');
end#
delimiter ;


delimiter #

create trigger admin_demotion
    after delete
    on reals.admins
    for each row
begin
    insert into notifications(user_id, text) values (old.user_id, 'You have been demoted to user');
end#
delimiter ;