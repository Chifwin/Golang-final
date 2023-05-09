call add_user('user', '1', 'Oleg', 'admin');
call add_user('seller', '1', 'Zarina', 'seller');
call add_user('buyer', '1', 'Oleg second', 'buyer');

insert into products(name, description)
values ('best prod', '');

insert into product_seller(product_id, seller_id, quantity, cost, published)
values ((select id from products limit 1),
        (select id from users where role = 'seller' limit 1),
        500, 10, true);

call buy((select id from users where role = 'buyer' limit 1), (select id from users where role = 'seller' limit 1), (select id from products limit 1), 5);

insert into scores values ((select id from purchases limit 1), 5, '');