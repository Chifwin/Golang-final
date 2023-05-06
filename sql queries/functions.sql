-- Authorize section
drop function if exists authorise_user(_username varchar, _password text);
create
    or replace function authorise_user(in in_username varchar, in in_password text)
    returns table
            (
                _id  integer,
                _username varchar,
                _name varchar,
                _role user_role
            )
    language plpgsql
as
$$
begin
    return query
        select id, username, name, role
        from users
        where username = in_username
          and password = crypt(in_password, password);
end;
$$;

drop procedure if exists add_user(_username varchar, _password text, _name varchar, _role user_role);
create
    or replace procedure add_user(in _username varchar, in _password text, _name varchar, role user_role)
    language plpgsql
as
$$
begin
    insert into users(username, password, name, role)
    values (_username, crypt(_password, gen_salt('bf')), _name, role);
end;
$$;