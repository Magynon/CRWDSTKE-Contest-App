-- Example create
CREATE TABLE IF NOT EXISTS products(
    id varchar(64) primary key,
    name varchar(64) not null,
    manufacturer varchar(64),
    price integer,
    stock integer,
    tags varchar(64)[]
    );

-- Example insert
insert into products(id, name, tags) values('sha256', 'Lapte', ARRAY['lactate', 'uht']);

-- Example selects
select * from products;
select id, name, tags from products ;
select id, name, from products where name='Lapte';

-- Example update
update products set price = 15000 WHERE name = 'Lapte';

-- Example delete
delete from products where name = 'Lapte';
