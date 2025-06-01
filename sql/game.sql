Create Database IF NOT EXISTS game;

DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '玩家ID',
  `server_no` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '服务器编号',
  `acct_name` varchar(25) NOT NULL DEFAULT '' COMMENT '账号名字',
  `lord_name` varchar(25) NOT NULL DEFAULT '' COMMENT '玩家领主名',
  `sign` varchar(100) NOT NULL DEFAULT '' COMMENT '签名',
  `birth_point` int(11) NOT NULL DEFAULT '0' COMMENT '出生点',
  `country` int(11) NOT NULL DEFAULT '0' COMMENT '所属国家',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '等级',
  `exp` int(11) NOT NULL DEFAULT '0' COMMENT '经验',
  `wood` int(11) NOT NULL DEFAULT '0' COMMENT '木材',
  `iron` int(11) NOT NULL DEFAULT '0' COMMENT '铁矿',
  `stone` int(11) NOT NULL DEFAULT '0' COMMENT '石料',
  `forage` int(11) NOT NULL DEFAULT '0' COMMENT '粮草',
  `gold` int(11) NOT NULL DEFAULT '0' COMMENT '金币',
  `diamond` int(11) NOT NULL DEFAULT '0' COMMENT '钻石',
  `bind_diamond` int(11) NOT NULL DEFAULT '0' COMMENT '绑定钻石',
  `decree` int(11) NOT NULL DEFAULT '0' COMMENT '政令',
  `army_order` int(11) NOT NULL DEFAULT '0' COMMENT '军令',
  `power` int(11) NOT NULL DEFAULT '0' COMMENT '势力值',
  `domain` int(11) NOT NULL DEFAULT '0' COMMENT '领地个数',
  `renown` int(11) NOT NULL DEFAULT '0' COMMENT '名望值',
  `renown_limit` int(11) NOT NULL DEFAULT '0' COMMENT '名望值上限',
  `feats` int(11) NOT NULL DEFAULT '0' COMMENT '武勋',
  `lord_level` int(11) NOT NULL DEFAULT '0' COMMENT '领主的爵位',
  `login_time` int(11) NOT NULL DEFAULT '0' COMMENT '登录时间戳',
  `logout_time` int(11) NOT NULL DEFAULT '0' COMMENT '退出时间戳',
  `register_time` int(11) NOT NULL DEFAULT '0' COMMENT '注册时间戳',
  `ally_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '同盟id',
  `ally_pos` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '同盟职位',
  `ally_name` varchar(25) NOT NULL DEFAULT '' COMMENT '同盟名字',
  `markLands`  text COMMENT '标记地块',
  `energy` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '体力精力',
  `normal_section` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '普通战役关卡id',
  `special_section` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '精英战役关卡id',
  `sectionMapStr`  text COMMENT '战役关卡次数',
  `sectionHeroStr`  text COMMENT '战役英雄列表',
  `lotteryStr`  text COMMENT '英雄抽卡',
  `first_occupy` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '玩家第一次占领等级土地',
  `soldierMapStr` text COMMENT '玩家兵种列表',
  `newGuideMapStr` text COMMENT '新手引导列表',
  `magicMapStr` text COMMENT '玩家魔法阵信息',
  `ally_request_cd` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '同盟申请CD',
  `recruit_cost` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '招募消耗钻石数',
  `ad_count` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '广告招募次数',
  `guide_reward_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '新手引导奖励id',
  `investigate_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '侦查时间戳',
  `tax_info`  text COMMENT '税收信息',
  `conscript_speed` int(11) NOT NULL DEFAULT '0' COMMENT '征兵加速',
  `conscript_speed_limit` int(11) NOT NULL DEFAULT '0' COMMENT '征兵加速上限',
  `speed_time_stamp` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '加速征兵时间戳',
  `recruit_ad_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '招募广告时间戳',
  `payment` int(11) DEFAULT '0' COMMENT '充值金额',
  `decree_time_stamp` int(11) DEFAULT '0' COMMENT '政令恢复时间戳',
  `resource_time_stamp` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '资源恢复时间戳',
  `zero_time_stamp` int(11) DEFAULT '0' COMMENT '0点时间戳',
  `four_time_stamp` int(11) DEFAULT '0' COMMENT '4点时间戳',
  `fail_compensate_max` text COMMENT '最大补偿次数',
  `common_flag` int(11) DEFAULT '0' COMMENT '通用的位存储标识',
  `black_market_reset_num` smallint(3) NOT NULL DEFAULT '0' COMMENT '黑市刷新次数',
  `language` smallint(3) NOT NULL DEFAULT '0' COMMENT '选择文本语言',
  `cooperation_info` text COMMENT '协助信息',
  `rename_times` smallint(3) NOT NULL DEFAULT '0' COMMENT '改名次数',
  `city_master_index` int(11) DEFAULT '0' COMMENT '城主府地块ID',
  `head_id` int(11) DEFAULT '0' COMMENT '头像id',
  `accumulate_payment` int(11) DEFAULT '0' COMMENT '累计充值金额（有些充值不算累充）',
  `section_small_box` text COMMENT '战役关卡小箱子未领取列表',
  `section_big_box` text COMMENT '战役关卡大箱子未领取列表',
  `section_power` smallint(3) NOT NULL DEFAULT '0' COMMENT '关卡能量值',
  `last_chest_stamp` int(11) NOT NULL DEFAULT '0' COMMENT '上次同盟宝箱领取时间戳',
  `contribute_sum` int(11) NOT NULL DEFAULT '0' COMMENT '同盟的个人贡献值',
  `last_daily_task_stamp` int(11) NOT NULL DEFAULT '0' COMMENT '刷出同盟个人任务的时间戳',
  `firstMapStr` text COMMENT '玩家的第一次',
  `villageStr` text COMMENT '玩家的村庄信息',
  `got_resident_count` int(11) NOT NULL DEFAULT '0' COMMENT '获得过的居民数量',
  `map_explore_daily_times` smallint(3) NOT NULL DEFAULT '0' COMMENT '大地图每日探索次数',
  `next_map_explore_stamp` int(11) NOT NULL DEFAULT '0' COMMENT '下次大地图探索时间戳',
  `begin_map_dispel_stamp` int(11) NOT NULL DEFAULT '0' COMMENT '开始大地图驱散时间戳',
  `map_dispel_land_id` int(11) NOT NULL DEFAULT '0' COMMENT '在驱散的地块id',
  `map_dispel_times` int(11) NOT NULL DEFAULT '0' COMMENT '每日已经驱散的次数',
  `map_dispel_land_list` text COMMENT '已经驱散的地块列表',
  `next_reset_hero_stamp` int(11) NOT NULL DEFAULT '0' COMMENT '下次重置英雄的时间戳',
  `first_occupy_award` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '玩家第一次占领等级土地奖励',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;


DROP TABLE IF EXISTS `t_user_city`;
CREATE TABLE `t_user_city` (
  `user_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '玩家ID',
  `build_type` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '建筑id',
  `level_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '建筑等级id',
  `resident_attr_type` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '建筑居民属性类型',
  `resident_attr_value` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '建筑居民属性',
  `produce_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '产出开始时间戳',
  `other_data` text COMMENT '研究院信息',
  PRIMARY KEY (`user_id`,`build_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `t_user_mail`;
CREATE TABLE `t_user_mail` (
  `mail_id` bigint(14) unsigned NOT NULL AUTO_INCREMENT COMMENT '邮件id',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '对应的玩家ID',
  `read_flag` smallint(3) NOT NULL DEFAULT '0' COMMENT '读取标志位，0未读，1已读',
  `get_flag` smallint(3) NOT NULL DEFAULT '0' COMMENT '领取标志位，0未领取，1已领取',
  `type` smallint(3) NOT NULL DEFAULT '0' COMMENT '邮件类型',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建邮件时间戳',
  `config_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '配置表id',
  `title` text COMMENT '邮件标题',
  `context` text COMMENT '邮件内容',
  `from_name` varchar(25) NOT NULL DEFAULT '' COMMENT '发件人名字',
  `params` text COMMENT '参数列表',
  `data_info` text COMMENT '附件',
  `url` text COMMENT '跳转url',
  `from_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '发件人的id，不是游戏玩家为0',
  PRIMARY KEY (`mail_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `t_war_mail`;
CREATE TABLE `t_war_mail` (
  `mail_id` bigint(14) unsigned NOT NULL AUTO_INCREMENT COMMENT '邮件id',
  `mail_type` smallint(3) NOT NULL DEFAULT '0' COMMENT '邮件类型,0为个人战报，1为同盟战报',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '对应的玩家ID',
  `create_time` bigint(14) NOT NULL DEFAULT '0' COMMENT '创建邮件时间戳',
  `store_flag` smallint(3) NOT NULL DEFAULT '0' COMMENT '0为正常，1为收藏',
  `land_id` int(11) NOT NULL DEFAULT '0' COMMENT '地块id',
  `war_type` smallint(3) NOT NULL DEFAULT '0' COMMENT ' 0防守，1进攻，2战役',
  `pk_flag` smallint(3) NOT NULL DEFAULT '0' COMMENT ' 是否交战（0为否，1为是）',
  `result` smallint(3) NOT NULL DEFAULT '0' COMMENT ' 胜负结果(0为成功，1为失败，2为平局，3为未战)',
  `occupy` smallint(3) NOT NULL DEFAULT '0' COMMENT ' 0为没占领，1为已经占领，2为失去占领，3成功沦陷，4被沦陷',
  `mail_series` bigint(14) NOT NULL DEFAULT '0' COMMENT '上一封邮件id（用于合并战报）',
  `read_flag` smallint(3) NOT NULL DEFAULT '0' COMMENT '读取标志位，0未读，1已读',
  `user_ally_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '对应的玩家同盟ID',
  `defend_refer_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '防守方地块refer_id',
  `attack_refer_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '攻击方地块refer_id',
  `ui_refer_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '防守地块显示refer_id',
  `land_name` varchar(25) NOT NULL DEFAULT '' COMMENT '地块名字',
  `a_war_player` text COMMENT 'A战斗玩家信息(包括展示的英雄)',
  `b_war_player` text COMMENT 'B战斗玩家信息(包括展示的英雄)',
  `detail_info` text COMMENT '更详细战报信息',
  `battle_time` int(11) NOT NULL DEFAULT '0' COMMENT '战斗耗时(毫秒)',
  PRIMARY KEY (`mail_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `t_store_war_mail`;
CREATE TABLE `t_store_war_mail` (
  `mail_id` bigint(14) unsigned NOT NULL AUTO_INCREMENT COMMENT '邮件id',
  `mail_type` smallint(3) NOT NULL DEFAULT '0' COMMENT '邮件类型,0为个人战报，1为同盟战报',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '对应的玩家ID',
  `create_time` bigint(14) NOT NULL DEFAULT '0' COMMENT '创建邮件时间戳',
  `store_flag` smallint(3) NOT NULL DEFAULT '0' COMMENT '0为正常，1为收藏',
  `land_id` int(11) NOT NULL DEFAULT '0' COMMENT '地块id',
  `war_type` smallint(3) NOT NULL DEFAULT '0' COMMENT ' 0防守，1进攻，2战役',
  `pk_flag` smallint(3) NOT NULL DEFAULT '0' COMMENT ' 是否交战（0为否，1为是）',
  `result` smallint(3) NOT NULL DEFAULT '0' COMMENT ' 胜负结果(0为成功，1为失败，2为平局，3为未战)',
  `occupy` smallint(3) NOT NULL DEFAULT '0' COMMENT ' 0为没占领，1为已经占领，2为失去占领，3成功沦陷，4被沦陷',
  `mail_series` bigint(14) NOT NULL DEFAULT '0' COMMENT '上一封邮件id（用于合并战报）',
  `read_flag` smallint(3) NOT NULL DEFAULT '0' COMMENT '读取标志位，0未读，1已读',
  `user_ally_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '对应的玩家同盟ID',
  `defend_refer_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '防守方地块refer_id',
  `attack_refer_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '攻击方地块refer_id',
  `ui_refer_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '防守地块显示refer_id',
  `land_name` varchar(25) NOT NULL DEFAULT '' COMMENT '地块名字',
  `a_war_player` text COMMENT 'A战斗玩家信息(包括展示的英雄)',
  `b_war_player` text COMMENT 'B战斗玩家信息(包括展示的英雄)',
  `detail_info` text COMMENT '更详细战报信息',
  `battle_time` int(11) NOT NULL DEFAULT '0' COMMENT '战斗耗时(毫秒)',
  PRIMARY KEY (`mail_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `t_user_task`;
CREATE TABLE `t_user_task` (
  `user_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '玩家ID',
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `type` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务类型',
  `subtype` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '子类型',
  `target` text NOT NULL COMMENT '任务目标',
  `task_state` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务状态',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '接任务时间戳',
  `expire_time` int(11) NOT NULL DEFAULT '0' COMMENT '有效(结束)时间戳',
  PRIMARY KEY (`user_id`,`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 物品表
-- -----------------------------------
DROP TABLE IF EXISTS `t_user_goods`;
CREATE TABLE `t_user_goods` (
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '对应的玩家ID',
  `goods_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '物品模版ID',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '类型',
  `subtype` int(11) NOT NULL DEFAULT '0' COMMENT '子类型',
  `num` int(11) NOT NULL DEFAULT '0' COMMENT '数量',
  `location` int(11) NOT NULL DEFAULT '0' COMMENT '位置1背包,2xx',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间戳',
  PRIMARY KEY (`user_id`,`goods_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 英雄表
-- -----------------------------------
DROP TABLE IF EXISTS `t_user_hero`;
CREATE TABLE `t_user_hero` (
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '对应的玩家ID',
  `hero_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄ID',
  `status`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄状态',
  `level`   int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄等级',
  `star`    int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄星级',
  `exp`     int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄经验',
  `grade`   int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄阶级',
  `grade_time`   int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄进阶完成时间戳',
  `skills`  text COMMENT '英雄技能',
  `equipments` text COMMENT '英雄装备',
  `strength`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄体力',
  `strength_time`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄体力更新时间戳',
  `soldier_num`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄兵力',
  `hurt_soldier`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄伤兵',
  `hurt_time`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄重伤时间戳',
  `army_id`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄所在部队的ID',
  `conscript_count`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄征兵的数量',
  `conscript_time`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '英雄征兵完成时间戳',
  PRIMARY KEY (`user_id`,`hero_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 部队表
-- -----------------------------------
DROP TABLE IF EXISTS `t_user_army`;
CREATE TABLE `t_user_army` (
  `user_id`     int(11) unsigned NOT NULL DEFAULT '0' COMMENT '对应的玩家ID',
  `army_id`     int(11) unsigned NOT NULL DEFAULT '0' COMMENT '部队ID',
  `hero_list`   text COMMENT '部队英雄列表',
  `max_hero_count`     int(11) unsigned NOT NULL DEFAULT '0' COMMENT '部队最大英雄数量',
  `army_name` varchar(25) NOT NULL DEFAULT '' COMMENT '部队名字',
  PRIMARY KEY (`user_id`,`army_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 同盟总表
-- -----------------------------------
DROP TABLE IF EXISTS `t_ally`;
CREATE TABLE `t_ally` (
  `ally_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '同盟id',
  `ally_icon` varchar(25) NOT NULL DEFAULT '' COMMENT '同盟icon',
  `ally_name` varchar(25) NOT NULL DEFAULT '' COMMENT '同盟名字',
  `ally_msg` text COMMENT '同盟描述',
  `leader_id` int(11) NOT NULL DEFAULT '0' COMMENT '盟主Id',
  `leader_name` varchar(25) NOT NULL DEFAULT '' COMMENT '盟主名字',
  `ally_level` int(11) NOT NULL DEFAULT '0' COMMENT '同盟等级',
  `ally_exp` int(11) NOT NULL DEFAULT '0' COMMENT '同盟经验',
  `ally_castlenum` int(11) NOT NULL DEFAULT '0' COMMENT '同盟城池数量',
  `country` int(11) NOT NULL DEFAULT '0' COMMENT '所属国家',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `ally_power` int(11) NOT NULL DEFAULT '0' COMMENT '同盟势力值',
  `change_time` int(11) NOT NULL DEFAULT '0' COMMENT '禅让截止时间戳',
  `change_userid` int(11) NOT NULL DEFAULT '0' COMMENT '禅让对象id',
  `country_attr`  text COMMENT '国家加成',
  `castle_attr`  text COMMENT '城池加成',
  `ally_relation`  text COMMENT '同盟外交',
  `fall_members`  text COMMENT '沦陷成员列表',
  `commands`  text COMMENT '指挥命令',
  `black_effect` int(11) NOT NULL DEFAULT '0' COMMENT '黑化影响',
  `break_up_time` int(11) NOT NULL DEFAULT '0' COMMENT '同盟解散时间',
  `cooperation_list` text COMMENT '协助列表',
  `invite_list` text COMMENT '同盟邀请列表',
  `last_release_task_stamp` int(11) NOT NULL DEFAULT '0' COMMENT '上次可发布任务恢复的时间戳',
  `release_task_times` int(11) NOT NULL DEFAULT '0' COMMENT '当前可发布次数',
  `ally_qa_info`  text COMMENT '同盟答题',
  `audit` smallint(3) NOT NULL DEFAULT '0' COMMENT '是否需要审核，0不需要审核，1需要审核',
  `treasure_hunt_map_stamp` int(11) NOT NULL DEFAULT '0' COMMENT '寻宝地图创建的时间戳',
  `treasure_hunt_map_id` int(11) NOT NULL DEFAULT '0' COMMENT '寻宝地图id',
  `ally_rename` smallint(3) NOT NULL DEFAULT '0' COMMENT '是否能改名，0不能改，1能改',
  PRIMARY KEY (`ally_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 同盟成员表
-- -----------------------------------
DROP TABLE IF EXISTS `t_ally_member`;
CREATE TABLE `t_ally_member` (
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '玩家',
  `ally_id` int(11) NOT NULL DEFAULT '0' COMMENT '同盟id',
  `contribute` int(11) NOT NULL DEFAULT '0' COMMENT '玩家贡献',
  `contribute_week` int(11) NOT NULL DEFAULT '0' COMMENT '一周玩家贡献',
  `feats_week` int(11) NOT NULL DEFAULT '0' COMMENT '一周玩家功勋',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `ally_pos` int(11) NOT NULL DEFAULT '0' COMMENT '职位：1普通会员，2管理员，3会长',
  `group_name` varchar(25) NOT NULL DEFAULT '' COMMENT '同盟分组名字',
  `group_leader` int(11) NOT NULL DEFAULT '0' COMMENT '同盟分组组长',
  `treasure_hunt_state` smallint(3) NOT NULL DEFAULT '0' COMMENT '同盟寻宝状态',
  `treasure_hunt_recover_stamp` int(11) NOT NULL DEFAULT '0' COMMENT '同盟寻宝开始恢复行动次数时间戳',
  `treasure_hunt_times` smallint(3) NOT NULL DEFAULT '0' COMMENT '同盟寻宝行动次数',
  `be_invite_hunt_times` smallint(3) NOT NULL DEFAULT '0' COMMENT '被邀请寻宝次数，日重置',
  `treasure_hunt_coordinate` int(11) NOT NULL DEFAULT '0' COMMENT '事发地坐标',
  `ally_qa_day_times` int(11) NOT NULL DEFAULT '0' COMMENT '同盟问答每日答题次数',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 同盟任务表
-- -----------------------------------
DROP TABLE IF EXISTS `t_ally_task`;
CREATE TABLE `t_ally_task` (
  `ally_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '同盟ID',
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `target` text NOT NULL COMMENT '任务目标',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '接任务时间戳',
  `duration` int(11) NOT NULL DEFAULT '0' COMMENT '持续时间',
  `member_list` text COMMENT '领取了的成员列表',
  PRIMARY KEY (`ally_id`,`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 玩家地块列表
-- -----------------------------------
DROP TABLE IF EXISTS `t_user_land`;
CREATE TABLE `t_user_land` (
  `land_id` int(11) NOT NULL DEFAULT '0' COMMENT '地块id',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '对应的玩家ID',
  PRIMARY KEY (`land_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- 创建玩家激活码表
-- ----------------------------
DROP TABLE IF EXISTS `t_exchange_code`;
CREATE TABLE `t_exchange_code` (
  `code` varchar(20) NOT NULL DEFAULT '' COMMENT '激活码',
  `platform` varchar(20) NOT NULL DEFAULT '' COMMENT '平台编号',
  `get_flag` int(11) NOT NULL DEFAULT '0' COMMENT '0未领取，1已领取',
  `gift_id` int(11) NOT NULL DEFAULT '0' COMMENT '礼包类型',
  `deal_time` int(11) DEFAULT '0' COMMENT '处理时间',
  `user_id` int(11) DEFAULT '0' COMMENT '领取人id',
  PRIMARY KEY (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 玩家活动
-- -----------------------------------
DROP TABLE IF EXISTS `t_user_activity`;
CREATE TABLE `t_user_activity` (
  `user_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '玩家ID',
  `sign_data` text COMMENT '签到活动',
  `sign_month` int(11) DEFAULT '0' COMMENT '签到月份',
  `online_gift` text COMMENT '在线活动',
  `common_flag` int(11) DEFAULT '0' COMMENT '通用活动领取标志位',
  `buy_goods_flag` int(11) DEFAULT '0' COMMENT '充值道具标志位',
  `purchase_limit` text COMMENT '限购道具列表',
  `month_card` text COMMENT '月卡',
  `accumulate_list` text COMMENT '累计充值领取列表',
  `fb_invite_list` text COMMENT 'fb领取列表',
  `black_market` text COMMENT '黑市',
  `power_guarantee_flag` smallint(3) NOT NULL DEFAULT '0' COMMENT '势力值低保0未领取，1已领取',
  `space_rift` text COMMENT '空间裂缝',
  `survey_list` text COMMENT '调查领取列表',
  `growth_fund` text COMMENT '成长基金领取状态',
  `recharge_daily_limit` text COMMENT '充值每日限购列表',
  `daily_gift` smallint(3) NOT NULL DEFAULT '0' COMMENT '每日礼包领取状态0未领取1领取',
  `rookie_count_limit` text COMMENT '新手期间使用次数',
  `happy_week` text COMMENT '七天乐活动信息',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 居民系统
-- -----------------------------------
DROP TABLE IF EXISTS `t_resident`;
CREATE TABLE `t_resident` (
  `resident_id` bigint(14) unsigned NOT NULL AUTO_INCREMENT COMMENT '居民id',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '对应的玩家ID',
  `gender` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '性别',
  `type` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '类型',
  `resident_attr`  text COMMENT '居民属性',
  `star` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '星级',
  `place` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '城内or城外',
  `boredom_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '厌倦结束时间戳',
  `birth_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '生育结束时间戳',
  `grow_up_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '成长结束时间戳',
  `lover_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '配偶',
  PRIMARY KEY (`resident_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 玩家反馈表
-- -----------------------------------
DROP TABLE IF EXISTS `t_feed_back`;
CREATE TABLE `t_feed_back` (
  `msg` text COMMENT '反馈内容',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '玩家Id',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 玩家封号表
-- -----------------------------------
DROP TABLE IF EXISTS `t_user_ban`;
CREATE TABLE `t_user_ban` (
   `user_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '玩家ID',
   `lord_name` varchar(25) NOT NULL DEFAULT '' COMMENT '玩家领主名',
   `ban_stamp` int(11) DEFAULT '0' COMMENT '封禁结束时间戳',
   `reason` text COMMENT '封禁原因',
   PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- -----------------------------------
-- 玩家充值表
-- -----------------------------------
DROP TABLE IF EXISTS `t_user_recharge`;
CREATE TABLE `t_user_recharge` (
  `index` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `order_no` varchar(32) NOT NULL DEFAULT '' COMMENT '账单Id',
  `user_id` int(11) NOT NULL COMMENT '玩家Id',
  `acct_name` varchar(25) NOT NULL DEFAULT '' COMMENT '账号名',
  `product_id` int(11) DEFAULT '0' COMMENT '充值的游戏商品id',
  `usd_amount` int(11) DEFAULT '0' COMMENT '用户支付的游戏道具以美元计价的金额，单位美元',
  `pay_stamp` int(11) DEFAULT '0' COMMENT '充值时间戳',
  PRIMARY KEY (`index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
