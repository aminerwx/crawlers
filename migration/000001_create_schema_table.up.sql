CREATE TABLE IF NOT EXISTS components (
                                  id serial primary key,
                                  name text,
                                  stock int,
                                  price int,
                                  current_price int,
                                  discount int,
                                  url text,
                                  menu text,
                                  menu_url text,
                                  category text,
                                  category_url text,
                                  subcategory text,
                                  subcategory_url text,
                                  platform text
);
