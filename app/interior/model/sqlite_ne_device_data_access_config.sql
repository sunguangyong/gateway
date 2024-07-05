CREATE TABLE ne_device_data_access_config
(
    nid                 INTEGER PRIMARY KEY AUTOINCREMENT,-- 主键ID，自增长，用作唯一标识符
    device_nid          INTEGER NOT NULL,-- 设备ID，不能为空
    config_type         INTEGER DEFAULT 0,-- 配置类型，整数，默认值为0，表示读/写权限，1表示只读权限，2表示只写权限
    config_id           TEXT,-- 配置ID 可以为空
    config_name         TEXT    NOT NULL,-- 配置名称，不为空
    endpoint            TEXT    NOT NULL,-- 连接URL，最大长度为3000字符，不为空
    protocol            TEXT    NOT NULL,-- 访问协议，最大长度为50字符，不为空。可能的值包括：ModbusTcp、ModbusRtu、OpcUa、HTTP、MTCONNECT、MITSUBISH_MC、SIEMENS_S7、OMRON_SINS
    json_access_options TEXT    NOT NULL,-- 访问连接基本配置，以JSON格式存储，不能为空
    timeout             INTEGER DEFAULT 1000,-- 连接超时时间，整数，默认值为1000
    agw_id              INTEGER,-- 所属网关ID，默认可以为空
    tenant_id           INTEGER,-- 租户ID，默认可以为空
    profile_nid         INTEGER,-- 设备配置文件ID，默认可以为空
    create_time         TEXT    NOT NULL,-- 创建时间，日期时间格式，不为空
    create_by           INTEGER NOT NULL,-- 创建者ID，不能为空
    issued              INTEGER DEFAULT 0,-- 是否已下发，布尔类型，默认值为0，表示未下发
    issue_time          TEXT,-- 下发时间，日期时间格式，可以为空
    UNIQUE (device_nid, config_name)                      -- 唯一键，由device_nid和config_name组成)
)