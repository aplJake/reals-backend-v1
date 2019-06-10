delimiter #

create trigger auto_create_profile after insert on reals.users
    for each row
    begin
        insert into reals.user_profile(user_id, profile_description) VALUES (new.user_id, '');
    end #
delimiter ;