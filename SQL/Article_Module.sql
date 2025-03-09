CREATE TABLE `article`  (
  `article_id` int NOT NULL AUTO_INCREMENT COMMENT '文章id',
  `article_title` varchar(255) NULL COMMENT '文章标题',
  `article_add_time` datetime NULL COMMENT '文章添加时间',
  `article_context` text NULL COMMENT '文章内容',
  `article_raise` int NULL COMMENT '文章点赞',
  `article_collection` int NULL COMMENT '文章收藏',
  `article_lookthrough` int NULL COMMENT '文章浏览次数',
  `author_id` int NULL COMMENT '【作者外键】',
  `article_status` tinyint NULL COMMENT '文章状态；0是草稿，1表示已发布，2表示已删除',
  `creat_time` datetime NULL COMMENT '创建时间',
  `update_time` datetime NULL COMMENT '最后更新时间',
  PRIMARY KEY (`article_id` DESC)
);

CREATE TABLE `article_collection`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `sys_user_id` int NULL COMMENT '【用户ID】',
  `article_id` int NULL COMMENT '【文章ID】',
  `collection_time` datetime NULL COMMENT '收藏时间',
  PRIMARY KEY (`id`)
);

CREATE TABLE `article_comment`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `article_id` int NULL COMMENT '【文章ID】',
  `user_id` int NULL COMMENT '【用户ID】',
  `parent_comment_id` int NULL COMMENT '父评论ID',
  `comment_context` text NULL COMMENT '评论内容',
  `comment_time` datetime NULL COMMENT '评论时间',
  `is_deleted` tinyint NULL COMMENT '是否删除（0正常；1删除）',
  PRIMARY KEY (`id`)
);

CREATE TABLE `article_like`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `sys_user_id` int NULL COMMENT '【用户ID】',
  `article_id` int NULL COMMENT '【文章ID】',
  `like_time` datetime NULL COMMENT '点赞时间',
  PRIMARY KEY (`id`)
);

CREATE TABLE `article_tag`  (
  `article_tag_id` int NOT NULL AUTO_INCREMENT COMMENT '文章标签id',
  `article_tag_name` varchar(255) NULL COMMENT '标签名称',
  `article_tag_addtime` datetime NULL COMMENT '标签添加时间',
  PRIMARY KEY (`article_tag_id`)
);

CREATE TABLE `article_tag_map`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `tag_id` int NULL COMMENT '【标签ID】',
  `article_id` int NULL COMMENT '【文章ID】',
  PRIMARY KEY (`id`)
);

CREATE TABLE `sys_user`  (
  `sys_user_id` int NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `sys_user_name` varchar(255) NULL COMMENT '用户名',
  `sys_user_passwd` varchar(255) NULL COMMENT '用户密码',
  `sys_user_register_time` datetime NULL COMMENT '注册时间',
  `sys_user_lastlogin` datetime NULL COMMENT '最后登录时间',
  `sys_user_avatar` varchar(255) NULL COMMENT '用户头像URL',
  `user_role` tinyint NULL COMMENT '用户角色权限',
  `creat_time` datetime NULL COMMENT '创建时间',
  `update_time` datetime NULL COMMENT '最后更新时间',
  PRIMARY KEY (`sys_user_id`)
);

CREATE VIEW `view_1` AS;

