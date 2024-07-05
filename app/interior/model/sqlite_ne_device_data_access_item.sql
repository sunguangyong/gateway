CREATE TABLE ne_device_data_access_item (
    nid INTEGER PRIMARY KEY AUTOINCREMENT, -- 主键
    device_nid          INTEGER NOT NULL,-- 设备ID，不能为空
    agw_id              INTEGER,-- 所属网关ID，默认可以为空
    tenant_id           INTEGER,-- 租户ID，默认可以为空
    config_nid INTEGER NOT NULL, -- ne_device_data_access_config.id
    config_type INTEGER NOT NULL DEFAULT 0, -- 0-读/写, 1-只读, 2-只写
    access_data TEXT NOT NULL, -- 读/写配置项(JSON对象)
    create_time TEXT NOT NULL -- 创建时间
);