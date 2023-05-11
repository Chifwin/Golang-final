-- Authorize section
drop function if exists authorise_user(_username varchar, _password text);
create
    or replace function authorise_user(in in_username varchar, in in_password text)
    returns table
            (
                _id       integer,
                _username varchar,
                _name     varchar,
                _role     user_role
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

-- Buy function
create or replace function buy(_buyer_id integer, _seller_id integer, _product_id integer, _quantity integer)
    returns table
            (
                like purchases
            )
    language plpgsql
as
$$
begin
    update product_seller set quantity = quantity - _quantity
    where product_id = _product_id
      and seller_id = _seller_id
      and published = true;
    if not FOUND then
        raise exception 'Seller do not have such product or it is not published';
        return;
    end if;

    return query insert into purchases (buyer_id, product_id, seller_id, quantity)
        values (_buyer_id, _product_id, _seller_id, _quantity) returning *;
end;
$$;