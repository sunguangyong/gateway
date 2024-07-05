CREATE TABLE `ne_device_data_access_item` (
    `nid` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `config_nid` bigint NOT NULL COMMENT 'ne_device_data_access_config.id',
    `device_nid` bigint NOT NULL COMMENT '设备ID',
    `agw_id` bigint NOT NULL COMMENT '所属网关ID',
    `tenant_id` bigint DEFAULT '0' COMMENT '租户ID',
    `config_type` int NOT NULL DEFAULT '0' COMMENT '0-读/写,1-只读,2-只写',
    `access_data` varchar(4000) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '读/写配置项(JSON对象)',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`nid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=148698 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC