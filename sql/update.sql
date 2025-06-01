use game;

-- -----------------------------------
-- 英雄表在字段grade后面增加字段grade_time
-- -----------------------------------
-- alter table t_user_hero add grade_time int(11) unsigned NOT NULL DEFAULT '0' after grade;

-- -----------------------------------
-- t_user表删除lord_sub_level字段
-- -----------------------------------
-- alter table t_user drop column lord_sub_level, drop column lord_task_list;
--
-- alter table t_user add magicMapStr text COMMENT '玩家魔法阵信息' after newGuideMapStr;

-- alter table t_war_mail add land_name text COMMENT '地块名字' after buser_ally_id;

-- 2019-10-31 数据库表更新
-- 玩家表t_user在字段fail_compense_max后面增加字段is_join_ally
-- 同盟表t_ally在字段black_effect后面增加字段break_up_time
-- 战报表t_war_mail和收藏战报表t_store_war_mail在字段attack_refer_id后面增加字段ui_refer_id
-- alter table t_user add is_join_ally smallint(3) NOT NULL DEFAULT '0' COMMENT ' 是否加入过同盟（0为否，1为是）' after fail_compense_max;
-- alter table t_ally add break_up_time int(11) NOT NULL DEFAULT '0' COMMENT '同盟解散时间' after black_effect;
-- alter table t_war_mail add ui_refer_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT '防守地块显示refer_id' after attack_refer_id;
-- alter table t_store_war_mail add ui_refer_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT '防守地块显示refer_id' after attack_refer_id;

-- 2019-11-5 数据库表更新
-- 同盟表t_user_activity在字段t_user_activity后面增加字段black_market
-- 同盟表t_user在字段is_join_ally后面增加字段black_market_reset_num
-- alter table t_user_activity add black_market text COMMENT '黑市' after fb_invite_list;
-- alter table t_user add black_market_reset_num smallint(3) NOT NULL DEFAULT '0' COMMENT '黑市刷新次数' after is_join_ally;

-- 2019-11-12 数据表更新
-- 玩家表t_user删除字段fallInfoStr沦陷信息
-- alter table t_user drop column fallInfoStr;

-- 2019-11-22 数据表更新
-- 玩家表t_user在字段black_market_reset_num后面增加字段language
-- 玩家表t_user_mail在字段data_info后面增加字段url
-- alter table t_user add language smallint(3) NOT NULL DEFAULT '0' COMMENT '选择文本语言' after black_market_reset_num;
-- alter table t_user_mail add url text COMMENT '跳转url' after data_info;

-- 2019-12-9 数据表更新
-- 玩家表t_user字段is_join_ally修改为字段common_flag
-- alter table t_user CHANGE is_join_ally common_flag int(11) DEFAULT '0' COMMENT '通用的位存储标识';

-- 2019-12-16 数据库表更新
-- 同盟表t_ally在字段break_up_time后面增加字段cooperation_list
-- 玩家表t_user在字段language后面增加字段cooperation_info
-- 玩家表t_user在字段cooperation_info后面增加字段rename_times
-- alter table t_ally add cooperation_list text COMMENT '协助列表' after break_up_time;
-- alter table t_user add cooperation_info text COMMENT '协助信息' after language;
-- alter table t_user add rename_times smallint(3) NOT NULL DEFAULT '0' COMMENT '改名次数' after cooperation_info;

-- 2020-2-6 数据库表更新
-- alter table t_user_activity add power_guarantee_flag smallint(3) NOT NULL DEFAULT '0' COMMENT '势力值低保0未领取，1已领取' after black_market;
-- alter table t_ally add invite_list text COMMENT '同盟邀请列表' after cooperation_list;

-- 2020-2-14 数据库表更新
-- alter table t_war_mail add battle_time int(11) NOT NULL DEFAULT '0' COMMENT '战斗耗时(毫秒)' after detail_info;
-- alter table t_store_war_mail add battle_time int(11) NOT NULL DEFAULT '0' COMMENT '战斗耗时(毫秒)' after detail_info;
-- alter table t_user_activity add space_rift text COMMENT '空间裂缝' after power_guarantee_flag;

-- 2020-01-19 数据库更新
-- 英雄表增加体力更新时间戳
-- alter table t_user_hero add strength_time int(11) NOT NULL DEFAULT '0' COMMENT '注册时间戳' after strength;

-- 2020-03-04 数据库更新
-- 玩家表增加头像id
-- alter table t_user add head_id int(11) DEFAULT '0' COMMENT '头像id' after city_master_index;

-- 2020-03-13 数据库更新
-- 邮件表增加发件人的id
-- alter table t_user_mail add from_id int(11) DEFAULT '0' COMMENT '发件人的id，不是游戏玩家为0' after url;

-- 2020-03-26 数据库更新
-- 部队表增加部队名称
-- alter table t_user_army add army_name varchar(25) NOT NULL DEFAULT '' COMMENT '部队名字' after max_hero_count;

-- 2020-04-17 数据库更新
-- -----------------------------------
-- 玩家充值表
-- -----------------------------------
-- # DROP TABLE IF EXISTS `t_user_recharge`;
-- # CREATE TABLE `t_user_recharge` (
-- #    `order_no` int(11) unsigned NOT NULL COMMENT '账单Id',
-- #    `user_id` int(11) NOT NULL COMMENT '玩家Id',
-- #    `lord_name` varchar(25) NOT NULL DEFAULT '' COMMENT '玩家领主名',
-- #    `product_id` int(11) DEFAULT '0' COMMENT '充值的游戏商品id',
-- #    `usd_amount` int(11) DEFAULT '0' COMMENT '用户支付的游戏道具以美元计价的金额，单位美元',
-- #    `pay_stamp` int(11) DEFAULT '0' COMMENT '充值时间戳',
-- #    PRIMARY KEY (`order_no`)
-- # ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- #
-- # alter table t_user add `accumulate_payment` int(11) DEFAULT '0' COMMENT '累计充值金额（有些充值不算累充）' after head_id;
--
-- -- 2020-04-25 数据库更新
-- alter table t_user_activity add survey_list text COMMENT '调查领取列表' after space_rift;

-- alter table t_user add `section_small_box` text COMMENT '战役关卡小箱子未领取列表' after accumulate_payment;
-- alter table t_user add `section_big_box` text COMMENT '战役关卡大箱子未领取列表' after section_small_box;

-- 2020-6-6 数据库更新
-- alter table t_user add last_chest_stamp int(11) NOT NULL DEFAULT '0' COMMENT '上次同盟宝箱领取时间戳' after section_power;
-- # alter table t_user add contribute_sum int(11) NOT NULL DEFAULT '0' COMMENT '同盟的个人贡献值' after last_chest_stamp;
-- # alter table t_ally add last_release_task_stamp int(11) NOT NULL DEFAULT '0' COMMENT '上次可发布任务的时间' after invite_list;
-- # DROP TABLE IF EXISTS `t_ally_task`;
-- # CREATE TABLE `t_ally_task` (
-- #   `ally_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '同盟ID',
-- #   `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
-- #   `target` text NOT NULL COMMENT '任务目标',
-- #   `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '接任务时间戳',
-- #   `duration` int(11) NOT NULL DEFAULT '0' COMMENT '持续时间',
-- #   `member_list` text COMMENT '领取了的成员列表',
-- #   PRIMARY KEY (`ally_id`,`task_id`)
-- # ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
--
-- -- 2020-6-15 数据库更新
-- # alter table t_user_activity add growth_fund text COMMENT '成长基金领取状态' after survey_list;
-- # alter table t_user_activity add recharge_daily_limit text COMMENT '成长基金领取状态' after growth_fund;
-- # alter table t_user_activity add daily_gift smallint(3) NOT NULL DEFAULT '0' COMMENT '每日礼包领取状态0未领取1领取' after recharge_daily_limit;
-- alter table t_ally add release_task_times int(11) NOT NULL DEFAULT '0' COMMENT '上次可发布任务的时间' after last_release_task_stamp;

-- 2020-6-16更新数据库
-- alter table t_ally add ally_qa_info  text COMMENT '同盟答题' after release_task_times;
-- alter table t_user add last_daily_task_stamp int(11) NOT NULL DEFAULT '0' COMMENT '同盟的个人贡献值' after contribute_sum;

-- 2020-6-29更新数据库
-- alter table t_user_activity add rookie_count_limit text COMMENT '新手期间使用次数' after daily_gift;
-- alter table t_ally add audit smallint(3) NOT NULL DEFAULT '0' COMMENT '是否需要审核，0不需要审核，1需要审核' after ally_qa_info;

-- 2020-7-15更新数据库
-- alter table t_ally add treasure_hunt_map_stamp int(11) NOT NULL DEFAULT '0' COMMENT '寻宝地图创建的时间戳' after audit;
-- alter table t_ally_member add treasure_hunt_state smallint(3) NOT NULL DEFAULT '0' COMMENT '同盟寻宝状态' after group_leader;
-- alter table t_ally_member add treasure_hunt_recover_stamp int(11) NOT NULL DEFAULT '0' COMMENT '同盟寻宝开始恢复行动次数时间戳' after treasure_hunt_state;
-- alter table t_ally_member add treasure_hunt_times smallint(3) NOT NULL DEFAULT '0' COMMENT '同盟寻宝行动次数' after treasure_hunt_recover_stamp;
-- alter table t_ally_member add be_invite_hunt_times smallint(3) NOT NULL DEFAULT '0' COMMENT '被邀请寻宝次数，日重置' after treasure_hunt_times;
-- alter table t_ally add treasure_hunt_map_id int(11) NOT NULL DEFAULT '0' COMMENT '寻宝地图id' after treasure_hunt_map_stamp;
-- alter table t_ally_member add treasure_hunt_coordinate int(11) NOT NULL DEFAULT '0' COMMENT '事发地坐标' after be_invite_hunt_times;

-- 2020-8-12更新数据库
-- alter table t_user_activity add happy_week_daily_gift smallint(3) NOT NULL DEFAULT '0' COMMENT '每日礼包领取状态0未领取1领取' after rookie_count_limit;

-- 2020-8-17更新数据库
-- alter table t_ally_member add ally_qa_day_times int(11) NOT NULL DEFAULT '0' COMMENT '同盟问答每日答题次数' after treasure_hunt_coordinate;

-- 2020-8-18更新数据库
-- alter table t_user add `firstMapStr` text COMMENT '玩家的第一次' after last_daily_task_stamp;

-- 2020-8-21更新数据库
-- alter table t_user_activity add happy_week_is_reward_mail smallint(3) NOT NULL DEFAULT '0' COMMENT '七天乐活动补发邮件是否已经发放0未发放1发放' after happy_week_daily_gift;

-- 2020-8-24更新数据库
-- alter table t_user add `villageStr` text COMMENT '玩家的村庄信息' after firstMapStr;
-- alter table t_user add `got_resident_count` text COMMENT '获得过的居民数量' after villageStr;

-- 2020-8-17更新数据库
-- alter table t_user add map_explore_daily_times smallint(3) NOT NULL DEFAULT '0' COMMENT '大地图每日探索次数' after got_resident_count;
-- alter table t_user add next_map_explore_stamp int(11) NOT NULL DEFAULT '0' COMMENT '下次大地图探索时间戳' after map_explore_daily_times;
-- alter table t_user add next_map_dispel_stamp int(11) NOT NULL DEFAULT '0' COMMENT '下次大地图驱散时间戳' after next_map_explore_stamp;

-- 2020-8-20更新数据库
-- alter table t_user add map_dispel_land_id int(11) NOT NULL DEFAULT '0' COMMENT '在驱散的地块id' after next_map_dispel_stamp;
-- alter table t_user add map_dispel_times int(11) NOT NULL DEFAULT '0' COMMENT '每日已经驱散的次数' after map_dispel_land_id;
-- alter table t_user add map_dispel_land_list text COMMENT '已经驱散的地块列表' after map_dispel_times;
-- alter table t_user add next_reset_hero_stamp int(11) NOT NULL DEFAULT '0' COMMENT '下次重置英雄的时间戳' after map_dispel_land_list;
-- alter table t_user add first_occupy_award int(11) unsigned NOT NULL DEFAULT '0' COMMENT '玩家第一次占领等级土地奖励' after next_reset_hero_stamp;

-- alter table t_user_activity add happy_week text COMMENT '七天乐活动信息' after rookie_count_limit;
# alter table t_user add `got_resident_count` text COMMENT '获得过的居民数量' after villageStr;
#
# -- 2020-8-17更新数据库
# alter table t_user add map_explore_daily_times smallint(3) NOT NULL DEFAULT '0' COMMENT '大地图每日探索次数' after got_resident_count;
# alter table t_user add next_map_explore_stamp int(11) NOT NULL DEFAULT '0' COMMENT '下次大地图探索时间戳' after map_explore_daily_times;
# alter table t_user add next_map_dispel_stamp int(11) NOT NULL DEFAULT '0' COMMENT '下次大地图驱散时间戳' after next_map_explore_stamp;
#
# -- 2020-8-20更新数据库
# alter table t_user add map_dispel_land_id int(11) NOT NULL DEFAULT '0' COMMENT '在驱散的地块id' after next_map_dispel_stamp;
# alter table t_user add map_dispel_times int(11) NOT NULL DEFAULT '0' COMMENT '每日已经驱散的次数' after map_dispel_land_id;
# alter table t_user add map_dispel_land_list text COMMENT '已经驱散的地块列表' after map_dispel_times;
# alter table t_user add next_reset_hero_stamp int(11) NOT NULL DEFAULT '0' COMMENT '下次重置英雄的时间戳' after map_dispel_land_list;
# alter table t_user add first_occupy_award int(11) unsigned NOT NULL DEFAULT '0' COMMENT '玩家第一次占领等级土地奖励' after next_reset_hero_stamp;

-- 2020-9-1更新数据库
# alter table t_ally add ally_rename smallint(3) NOT NULL DEFAULT '0' COMMENT '是否能改名，0不能改，1能改' after treasure_hunt_map_id;
# alter table t_user_activity drop column happy_week_daily_gift, drop column happy_week_is_reward_mail;
alter table t_user drop column recruit_free_time, drop column recruit_free_count;

# 2020-9-7更新数据库
alter table t_user_recharge CHANGE lord_name acct_name varchar(25) NOT NULL DEFAULT '' COMMENT '账号名';
