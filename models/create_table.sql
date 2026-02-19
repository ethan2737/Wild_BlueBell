-- 使用数据库
use wild_bluebell;

/*
 * 表名：user
 * 描述：用户基础信息表
 *
 * 设计说明：
 * 1. id 为数据库内部自增主键，仅用于表关联
 * 2. user_id 为业务侧用户唯一标识，对外暴露使用
 * 3. username 全局唯一，用于登录（不区分大小写）
 * 4. password 存储加密后的密码（不可明文）
 * 5. gender 使用枚举值表示性别
 * 6. create_time / update_time 由数据库自动维护
 * 7. collate 用户字符串的比较和排序
 * 8. tinyint 较小的整数，适合分类、枚举
 * 9. 将字段使用 `` 包裹，主要解决关键字冲突，提高工程可用性，如果不加也可正常执行，但会存在安全风险
 */
CREATE TABLE `user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键（数据库内部使用）',
    `user_id` bigint(20) NOT NULL COMMENT '业务用户ID（对外唯一标识）',
    `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '登录用户名（唯一，不区分大小写）',
    `password` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户密码（加密存储）',
    `email` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户邮箱',
    `gender` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别：0未知 1男 2女',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`), -- 主键不支持备注
    UNIQUE KEY `idx_username` (`username`) USING BTREE, -- 索引不支持备注
    UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户基础信息表';

-- 防止community表存在
drop table if exists `community`;

-- 创建community表
create table `community` (
    `id` int(11) not null auto_increment,
    `community_id` int(10) unsigned not null,
    `community_name` varchar(100) collate utf8mb4_general_ci not null,
    `introduction` varchar(256) collate utf8mb4_general_ci not null,
    `create_time` timestamp not null default current_timestamp,
    `update_time` timestamp not null default current_timestamp on update current_timestamp,
    primary key (`id`),
    unique key `idx_community_id` (`community_id`),
    unique key `idx_community_name` (`community_name`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

-- 插入默认数据
insert into `community` values ('1', '1', 'Go', 'Golang', '2025-02-13 05:22:36', '2025-02-13 05:22:36');
insert into `community` values ('2', '2', 'leetcode', '刷题刷题', '2025-02-13 05:22:36', '2025-02-13 05:22:36');
insert into `community` values ('3', '3', 'CS:GO', 'Rush B....', '2025-02-13 05:22:36', '2025-02-13 05:22:36');
insert into `community` values ('4', '4', 'LOL', '欢迎来到英雄联盟！', '2025-02-13 05:22:36', '2025-02-13 05:22:36');

-- 判断帖子表是否存在
drop table if exists `post`;
-- 创建帖子表
create table `post` (
    `id` bigint(20) not null auto_increment,
    `post_id` bigint(20) not null comment '帖子ID',
    `title` varchar(128) collate utf8mb4_general_ci not null  comment '标题',
    `content` varchar(8192) collate utf8mb4_general_ci not null comment '内容',
    `author_id` bigint(20) not null comment '作者的用户id',
    `community_id` bigint(20) not null comment '所属社区',
    `status` tinyint(4) not null default '1' comment '帖子状态',
    `create_time` timestamp null default current_timestamp comment '创建时间',
    `update_time` timestamp null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (`id`),
    unique key `idx_post_id` (`post_id`),
    key `idx_author_id` (`author_id`),
    key `idx_community_id` (`community_id`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;