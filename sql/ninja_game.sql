Create Database IF NOT EXISTS ninja_game;

DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '玩家ID',
  `acct_name` varchar(64) NOT NULL DEFAULT '' COMMENT '账号',
  `server_no` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '服务器编号',
  `sex` int(11) NOT NULL DEFAULT 0 COMMENT '性别',
  `username` varchar(25) NOT NULL DEFAULT '' COMMENT '玩家领主名',
  `name_index`varchar(25) NOT NULL DEFAULT '-' COMMENT '名称索引',
  `user_photo` smallint(4)  NOT NULL DEFAULT 0 COMMENT '头像',
  `user_border` smallint(4)  NOT NULL DEFAULT 0 COMMENT '相框',
  `gold` int(11) NOT NULL DEFAULT '0' COMMENT '金币',
  `gemstone` int(11) NOT NULL DEFAULT '0' COMMENT '宝石',
  `level` int(11) NOT NULL DEFAULT 0 COMMENT '等级',
  `exp` int(11) NOT NULL DEFAULT 0 COMMENT '经验',
  `max_exp` int(11) NOT NULL DEFAULT 0 COMMENT '最大经验',
  `rank_id` smallint(4) NOT NULL DEFAULT 1 COMMENT '当前段位',
  `star` smallint(4) NOT NULL DEFAULT 0 COMMENT '星星',
  `star_count` smallint(4) NOT NULL DEFAULT 0 COMMENT '总星级 每一段位，每一星都能对应一个数',
  `his_rank_id` smallint(4)  NOT NULL DEFAULT 0 COMMENT '历史最高段位',
  `ninja_id` smallint(4) NOT NULL DEFAULT 1 COMMENT '忍阶',
  `ninja_id_gift` smallint(4) NOT NULL DEFAULT 0 COMMENT '忍阶礼包',
  `archive_point` int(11) NOT NULL DEFAULT 0 COMMENT '成就点',
  `max_archive_point` int(11) NOT NULL DEFAULT 0 COMMENT '最大成就点',
  `use_skin` smallint(4) NOT NULL DEFAULT 1 COMMENT '使用的皮肤',
  `got_skins` char(64) NOT NULL DEFAULT '1' COMMENT '皮肤ids:2,3',
  `game_duration` int(11) NOT NULL DEFAULT 0 COMMENT '游戏时长',
  `match_game_num` int(11)  NOT NULL DEFAULT 0 COMMENT '排位总场数',
  `match_win_num` int(11)  NOT NULL DEFAULT 0 COMMENT '排位胜场数',
  `match_wolf_num`int(11)  NOT NULL DEFAULT 0 COMMENT '狼人总场数',
  `wolf_win_num` int(11)  NOT NULL DEFAULT 0 COMMENT '狼人胜场数',
  `poor_win_num` int(11)  NOT NULL DEFAULT 0 COMMENT '平民胜场数',
  `offline_num` int(11)  NOT NULL DEFAULT 0 COMMENT '掉线局数',
  `vote_total` int(11)  NOT NULL DEFAULT 0 COMMENT '总投票次数',
  `vote_correct_total` int(11)  NOT NULL DEFAULT 0 COMMENT '投票正确次数',
  `vote_failed_total` int(11)  NOT NULL DEFAULT 0 COMMENT '投票错误次数',
  `kill_total` int(11)  NOT NULL DEFAULT 0 COMMENT '杀人次数',
  `wolf_kill_total` int(11)  NOT NULL DEFAULT 0 COMMENT '狼人杀人次数',
  `pool_kill_total` int(11)  NOT NULL DEFAULT 0 COMMENT '平民杀人次数',
  `bekilled_total` int(11)  NOT NULL DEFAULT 0 COMMENT '被杀害次数',
  `bevoteed_total` int(11)  NOT NULL DEFAULT 0 COMMENT '被票杀次数',
  `update_time` int(11) NOT NULL DEFAULT 0 COMMENT '最后更新时间',
  `room_id` int(11) NOT NULL DEFAULT 0 COMMENT '房间id',
  `openid` varchar(255) NOT NULL DEFAULT '' COMMENT '渠道用户唯一标识',
  `channel` varchar(16) NOT NULL DEFAULT '' COMMENT '渠道',
  `sex_modify` int(11) NOT NULL DEFAULT '1' COMMENT '性别修改',
  `fresh_gift_step` smallint(4) NOT NULL DEFAULT 0 COMMENT '新手礼包进度0-未达成，1-倒计时，2-已完成',
  `fresh_end_time` int(11) NOT NULL DEFAULT 0 COMMENT '倒计时结束时间',
  `register_time` int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
   PRIMARY KEY (`user_id`),
   INDEX idx_star_count(`star_count`);
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `t_user_title`;
CREATE TABLE `t_user_title` (
  `user_id` int(11) NOT NULL COMMENT '玩家ID',
  `keep_first_out`int(11) NOT NULL DEFAULT 0 COMMENT '连续第一个出局数',
  `keep_wolf` int(11) NOT NULL DEFAULT 0 COMMENT '连续狼人数',
  `keep_poor` int(11) NOT NULL DEFAULT 0 COMMENT '连续平民数',
  `keep_noitem` int(11) NOT NULL DEFAULT 0 COMMENT '连续未得到道具',
  `total_kill_poor` int(11) NOT NULL DEFAULT 0 COMMENT '累计杀平民',
  `total_wolf_day` int(11) NOT NULL DEFAULT 0 COMMENT '累计几天狼人',
  `wolf_timestamp` int(11) NOT NULL DEFAULT 0 COMMENT '累计狼人日期',
  `total_task` int(11) NOT NULL DEFAULT 0 COMMENT '任务次数',
  `total_soul_task` int(11) NOT NULL DEFAULT 0 COMMENT '灵魂状态帮助队友次数',
  `total_gold` int(11) NOT NULL DEFAULT 0 COMMENT '总金币',
  `total_ad` int(11) NOT NULL DEFAULT 0 COMMENT '广告次数',
  `total_archive` int(11) NOT NULL DEFAULT 0 COMMENT '总成就点',
  `use_title` int(11) NOT NULL DEFAULT 0 COMMENT '使用的称号',
  `got_titles` varchar(256) NOT NULL DEFAULT '' COMMENT '获得的称号，称号id:201,201',
  `skin_red_datas` varchar(256) NOT NULL DEFAULT '' COMMENT '皮肤红点数据 1,2',
  `title_red_datas` varchar(256) NOT NULL DEFAULT '' COMMENT '称号红点数据 1,2',
   PRIMARY KEY (`user_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `t_user_archive`;
CREATE TABLE `t_user_archive` (
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '玩家ID',
  `archive_type` smallint(4) NOT NULL DEFAULT 0 COMMENT '成就类型',
  `archive_id` smallint(4) NOT NULL DEFAULT 0 COMMENT '成就项目',
  `archive_next` smallint(4) NOT NULL DEFAULT 0 COMMENT '下一项达成',
  `got_status` smallint(4) NOT NULL DEFAULT 0 COMMENT '达成状态: 0-已领取,1-未领取',
   PRIMARY KEY (`user_id`, `archive_type`, `archive_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `t_item`;
CREATE TABLE `t_item` (
    `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
    `item_id` int(11) NOT NULL DEFAULT '0' COMMENT '道具id',
    `num` int(11) NOT NULL DEFAULT '0' COMMENT '道具数量',
    `user_id` int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`uid`),
    KEY `idx_user_id` (`user_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `t_user_mail`;
CREATE TABLE `t_user_mail` (
    `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
    `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '用户id',
    `detail_type` smallint(4) NOT NULL DEFAULT 0 COMMENT '详情类型，1-key类型需转义',
    `title` varchar(64) NOT NULL DEFAULT '' COMMENT '标题',
    `content` varchar(256) NOT NULL DEFAULT '' COMMENT '内容',
    `annex_items` varchar(256) NOT NULL DEFAULT 0 COMMENT '附件',
    `annex_open` smallint(4) NOT NULL DEFAULT 0 COMMENT '附件领取 0-否 1-打開',
    `is_read` smallint(4) NOT NULL DEFAULT 0 COMMENT '是否已读：0-否，1-已过期',
    `is_expire` smallint(4) NOT NULL DEFAULT 0 COMMENT '是否过期：0-否，1-已过期',
    `create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `update_time` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
    PRIMARY KEY (`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `t_store`;
CREATE TABLE `t_store`(
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

DROP TABLE IF EXISTS `t_acct`;
CREATE TABLE `t_acct`(
    `acct_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '账号id',
    `acct_name` varchar(64) NOT NULL DEFAULT '' COMMENT '账号',
    `password` varchar(256) NOT NULL DEFAULT '',
    PRIMARY KEY (`acct_id`),
    UNIQUE index_acct_name (acct_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

DROP TABLE IF EXISTS `t_world_chat`;
CREATE TABLE `t_world_chat`(
    `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
    `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '用户id',
    `username` varchar(25) NOT NULL DEFAULT '' COMMENT '玩家领主名',
    `user_photo` smallint(4)  NOT NULL DEFAULT 0 COMMENT '相框',
    `chat_data` varchar(128) NOT NULL DEFAULT '' COMMENT '聊天内容',
    `create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

DROP TABLE IF EXISTS `t_guide`;
CREATE TABLE `t_guide` (
   `user_id` int(11) NOT NULL DEFAULT 0,
   `guide_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '新手引导串',
   `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
   PRIMARY KEY (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

DROP TABLE IF EXISTS `t_action`;
CREATE TABLE `t_action` (
   `user_id` int(11) NOT NULL DEFAULT 0,
   `actions` varchar(500) NOT NULL DEFAULT '' COMMENT '行为',
   `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
   PRIMARY KEY (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

DROP TABLE IF EXISTS `t_dan_stat`;
CREATE TABLE `t_dan_stat` (
  `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
  `rank_id` int(11) NOT NULL DEFAULT '0' COMMENT '段位',
  `star` int(11) NOT NULL DEFAULT '0' COMMENT '星',
  `register_date` varchar(50) NOT NULL DEFAULT '' COMMENT '注册日期',
  PRIMARY KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

DROP TABLE IF EXISTS `t_game_data`;
CREATE TABLE `t_game_data` (
   `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
   `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
   `username` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
   `game_num` int(11) NOT NULL DEFAULT '0' COMMENT '游戏局数',
   `login_date` varchar(50) NOT NULL DEFAULT '' COMMENT '登录日期',
   `register_time` varchar(50) NOT NULL DEFAULT '0' COMMENT '注册时间',
   `update_time` varchar(50) NOT NULL DEFAULT '0' COMMENT '更新时间',
   PRIMARY KEY (`uid`),
   KEY `idx_user_id` (`user_id`),
   KEY `idx_login_date` (`login_date`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

DROP TABLE IF EXISTS `t_login_data`;
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

DROP TABLE IF EXISTS `t_ad_data`;
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

DROP TABLE IF EXISTS `t_win_data`;
CREATE TABLE `t_win_data` (
      `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
      `wolf_man` int(11) NOT NULL DEFAULT '0' COMMENT '狼人赢',
      `normal_man` int(11) NOT NULL DEFAULT '0' COMMENT '平民赢',
      `stat_date` varchar(50) NOT NULL DEFAULT '' COMMENT '统计日期',
      PRIMARY KEY (`uid`),
      KEY `idx_stat_date` (`stat_date`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

DROP TABLE IF EXISTS `t_kill_data`;
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

DROP TABLE IF EXISTS `t_confer_data`;
CREATE TABLE `t_confer_data` (
     `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
     `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
     `rank_id` int(11) NOT NULL DEFAULT '0' COMMENT '段位',
     `confer_num` int(11) NOT NULL DEFAULT '0' COMMENT '会议次数',
     `stat_date` varchar(50) NOT NULL DEFAULT '' COMMENT '统计日期',
     `register_time` int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
     PRIMARY KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;

DROP TABLE IF EXISTS `t_item_data`;
CREATE TABLE `t_item_data` (
       `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增序号',
       `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
       `item_id` int(11) NOT NULL DEFAULT '0' COMMENT '道具id',
       `num` int(11) NOT NULL DEFAULT '0' COMMENT '数量',
       `add_time` int(11) NOT NULL DEFAULT '0' COMMENT '添加时间',
       PRIMARY KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 ROW_FORMAT = COMPACT;