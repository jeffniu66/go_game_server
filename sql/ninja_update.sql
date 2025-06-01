-- 20201125 jeffrey
alter table t_user add `lord_name` varchar(25) NOT NULL DEFAULT '' COMMENT '玩家领主名';
alter table t_user add `gold` int(11) NOT NULL DEFAULT '0' COMMENT '金币';

alter table t_user add `rank` int(11) NOT NULL DEFAULT 1 COMMENT '段位';
alter table t_user add `stage` int(11) NOT NULL DEFAULT 1 COMMENT '阶段';
alter table t_user add `star` int(11) NOT NULL DEFAULT 0 COMMENT '星星';
alter table t_user add `star_count` int(11) NOT NULL DEFAULT 0 COMMENT '总星级 每一段位，每一星都能对应一个数';
alter table t_user add `level` int(11) NOT NULL DEFAULT 1 COMMENT '等级';
alter table t_user add `ninja_level` smallint(4) NOT NULL DEFAULT 1 COMMENT '忍阶';
alter table t_user add `ninja_star` smallint(4) NOT NULL DEFAULT 1 COMMENT '忍星';
alter table t_user add `ninja_point` smallint(4) NOT NULL DEFAULT 0 COMMENT '忍度点';
alter table t_user add `use_skin` smallint(4) NOT NULL DEFAULT 0 COMMENT '使用的皮肤';
alter table t_user add `got_skins` char(64) NOT NULL DEFAULT '' COMMENT '皮肤ids:2,3';
alter table t_user add `game_duration` int(11) NOT NULL DEFAULT 0 COMMENT '游戏时长';
alter table t_user add `match_game_num` smallint(4)  NOT NULL DEFAULT 0 COMMENT '排位总局数';
alter table t_user add `match_win_num` smallint(4)  NOT NULL DEFAULT 0 COMMENT '排位胜场局数';
alter table t_user add `vote_total` smallint(4)  NOT NULL DEFAULT 0 COMMENT '总投票次数';
alter table t_user add `vote_correct_total` smallint(4)  NOT NULL DEFAULT 0 COMMENT '投票正确次数';
alter table t_user add `kill_total` smallint(4)  NOT NULL DEFAULT 0 COMMENT '杀人次数';
alter table t_user add `bekilled_total` smallint(4)  NOT NULL DEFAULT 0 COMMENT '被杀害次数';
alter table t_user add `bevoteed_total` smallint(4)  NOT NULL DEFAULT 0 COMMENT '被杀票杀次数';
alter table t_user add `user_photo` varchar(64)  NOT NULL DEFAULT '' COMMENT '头像';
alter table t_user add `exp` int(11) NOT NULL DEFAULT 0 COMMENT '经验';
alter table t_user add `max_exp` int(11) NOT NULL DEFAULT 0 COMMENT '最大经验';

-- 20201209 jeffrey
CREATE TABLE `t_item` (
    `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
    `item_id` int(11) NOT NULL DEFAULT '0' COMMENT '道具id',
    `num` int(11) NOT NULL DEFAULT '0' COMMENT '道具数量',
    `user_id` int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- 20201216 jeffrey
alter table t_user add `gemstone` int(11) NOT NULL DEFAULT '0' COMMENT '宝石' after gold;

-- 20201219 jeffrey
CREATE TABLE `t_store` (
    `user_id` int(11) NOT NULL DEFAULT 0,
    `box_use_ad_num` int(11) NOT NULL DEFAULT '0',
    `box_last_use_ad_time` int(11) NOT NULL DEFAULT '0',
    `mys_skin_buy_num` int(11) NOT NULL DEFAULT '0',
    `mys_skin_last_refresh_time` int(11) NOT NULL DEFAULT '0',
    `mys_skin_chip_id` int(11) NOT NULL DEFAULT '0',
    `skin_use_ad_num` int(11) NOT NULL DEFAULT '0',
    `skin_last_refresh_time` int(11) NOT NULL DEFAULT '0',
    `gold_view_ad_num` int(11) NOT NULL DEFAULT '0',
    `gold_last_view_ad_time` int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210113 jeffrey
alter table t_user add `openid` varchar(255) NOT NULL DEFAULT '' COMMENT '渠道用户唯一标识';

-- 20210122
alter table t_user add `channel` varchar(16) NOT NULL DEFAULT '' COMMENT '渠道';

-- 20210126 jeffrey
CREATE TABLE `t_guide` (
    `user_id` int(11) NOT NULL DEFAULT 0,
    `guide_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '新手引导串',
    PRIMARY KEY (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210129 jeffrey
alter table `t_guide` add `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间';

-- 20210130 jeffrey
CREATE TABLE `t_action` (
    `user_id` int(11) NOT NULL DEFAULT 0,
    `actions` varchar(500) NOT NULL DEFAULT '' COMMENT '行为',
    `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
    PRIMARY KEY (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210201 jeffrey
alter table t_user add `sex_modify` int(11) NOT NULL DEFAULT '1' COMMENT '性别修改';

-- 20210201
alter table t_user add `fresh_gift_step` smallint(4) NOT NULL DEFAULT 0 COMMENT '新手礼包进度0-未达成，1-倒计时，2-已完成';
alter table t_user add `fresh_end_time` int(11) NOT NULL DEFAULT 0 COMMENT '倒计时结束时间';

-- 20210202 jeffrey
alter table t_user add `register_time` int(11) NOT NULL DEFAULT '0' COMMENT '注册时间';

-- 20210219 jeffrey
CREATE TABLE `t_dan_stat` (
    `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
    `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
    `username` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
    `rank_id` int(11) NOT NULL DEFAULT '0' COMMENT '段位',
    `star` int(11) NOT NULL DEFAULT '0' COMMENT '星',
    `register_date` varchar(50) NOT NULL DEFAULT '' COMMENT '注册日期',
    PRIMARY KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210222 jeffrey
CREATE TABLE `t_game_data` (
     `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
     `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
     `username` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
     `game_num` int(11) NOT NULL DEFAULT '0' COMMENT '游戏局数',
     `login_date` varchar(50) NOT NULL DEFAULT '' COMMENT '登录日期',
     `register_time` int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
     `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
     PRIMARY KEY (`uid`),
     KEY `idx_user_id` (`user_id`),
     KEY `idx_login_date` (`login_date`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210223 jeffrey
CREATE TABLE `t_login_data` (
   `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
   `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
   `username` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
   `login_date` varchar(50) NOT NULL DEFAULT '' COMMENT '登录日期',
   `register_time` int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
   PRIMARY KEY (`uid`),
   KEY `idx_user_id` (`user_id`),
   KEY `idx_login_date` (`login_date`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210223 jeffrey
CREATE TABLE `t_ad_data` (
    `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
    `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
    `ad_type` int(11) NOT NULL DEFAULT '0' COMMENT '广告类型',
    `ad_num` int(11) NOT NULL DEFAULT '0' COMMENT '广告数量',
    `register_time` int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
    `stat_date` varchar(50) NOT NULL DEFAULT '' COMMENT '统计日期',
    PRIMARY KEY (`uid`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_ad_type` (`ad_type`),
    KEY `idx_stat_date` (`stat_date`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210224 jeffrey
CREATE TABLE `t_win_data` (
    `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
    `wolf_man` int(11) NOT NULL DEFAULT '0' COMMENT '狼人赢',
    `normal_man` int(11) NOT NULL DEFAULT '0' COMMENT '平民赢',
    `stat_date` varchar(50) NOT NULL DEFAULT '' COMMENT '统计日期',
    PRIMARY KEY (`uid`),
    KEY `idx_stat_date` (`stat_date`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210224 jeffrey
CREATE TABLE `t_kill_data` (
      `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
      `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
      `rank_id` int(11) NOT NULL DEFAULT '0' COMMENT '段位',
      `kill_num` int(11) NOT NULL DEFAULT '0' COMMENT '杀人数',
      `stat_date` varchar(50) NOT NULL DEFAULT '' COMMENT '统计日期',
      PRIMARY KEY (`uid`),
      KEY `idx_user_id` (`user_id`),
      KEY `idx_stat_date` (`stat_date`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210224 jeffrey
CREATE TABLE `t_confer_data` (
      `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
      `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
      `rank_id` int(11) NOT NULL DEFAULT '0' COMMENT '段位',
      `confer_num` int(11) NOT NULL DEFAULT '0' COMMENT '会议次数',
      `stat_date` varchar(50) NOT NULL DEFAULT '' COMMENT '统计日期',
      `register_time` int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
      PRIMARY KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210226 jeffrey
CREATE TABLE `t_item_data` (
     `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
     `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
     `item_id` int(11) NOT NULL DEFAULT '0' COMMENT '道具id',
     `num` int(11) NOT NULL DEFAULT '0' COMMENT '数量',
     `add_time` int(11) NOT NULL DEFAULT '0' COMMENT '添加时间',
     PRIMARY KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

-- 20210302 jeffrey
ALTER TABLE `t_item` ADD INDEX idx_user_id(`user_id`);

-- 20210303
ALTER TABLE `t_user_title` ADD `skin_red_datas` varchar(256) NOT NULL DEFAULT '' COMMENT '皮肤红点数据 1,2';
ALTER TABLE `t_user_title` ADD `title_red_datas` varchar(256) NOT NULL DEFAULT '' COMMENT '称号红点数据 1,2';
ALTER TABLE `t_user` ADD INDEX idx_star_count(`star_count`);

------------------------------------------------v1.0.2 end-------------------------------------------
-----------------------------------------------------------------------------------------------------
