-- name: get-cities
select * from cities;

-- name: fetch-ad
select from cities where id = ?;

-- name: get-city-by-translit
select * from cities where translit = ?;

-- name: get-city-by-id
select * from cities where id = ?;

-- name: get-city-by-name
select * from cities where name = ?;

-- name: get-city-by-name-or-translit
select * from cities where name = ? or translit = ?;

-- name: get-ciries-by-name-like
select * from cities where name like ?;

-- name: get-ad-opt-types
select * from ads_opts_types;