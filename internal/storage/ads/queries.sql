-- name: get-ad-opt-types
select * from ads_opts_types;

-- name: fetch-ad
select user_id from ads where id = ?;

-- name: fetch-ads
select user_id from ads where id in (?);

-- name: fetch-page-ads
select id, user_id from ads limit ? offset ?;

-- name: fetch-ads-total-pages-count
select ceil(count(*) / ?) from ads;

-- name: fetch-ads-opts
select ad_id, opt_type, opt_value_str, opt_value_uint, opt_value_bool, opt_value_blob
from ads_opts
where ad_id in (?)
  and opt_type not in (?);

-- name: fetch-ad-string-opt
select opt_value_str
from ads_opts
where ad_id = ?
  and opt_type = ?;

-- name: get-recommended-ads
select id, user_id from ads order by rand() limit ?;

-- name: get-recommended-city-ads
select ad_id as id, 0 from ads_opts where opt_type = 58 and opt_value_str = ? order by rand() limit ?;