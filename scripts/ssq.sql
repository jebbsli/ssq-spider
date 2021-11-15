CREATE TABLE sys_config_t (
    type_id VARCHAR(64),
    name VARCHAR(64),
    value VARCHAR(128)
) DEFAULT CHARSET=utf8;

CREATE TABLE ssq_number_t (
    open_no VARCHAR(32),
    red_num VARCHAR(32),
    blue_num VARCHAR(4),
    ball_sort VARCHAR(32)
) DEFAULT CHARSET=utf8;