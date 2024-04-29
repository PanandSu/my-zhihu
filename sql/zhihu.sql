drop database if exists zhihu;
create database zhihu;
use zhihu;
#
drop table if exists users;
create table users
(
    id              int(11) unsigned not null auto_increment,
    gender          tinyint(2)       not null default '-1',
    avatar_url      varchar(100)     not null default '/static/images/default.jpg',
    email           varchar(50)      not null default '' unique,
    fullname        varchar(100)     not null,
    password        varchar(100)     not null,
    headline        varchar(500)              default '',
    url_token       varchar(50)      not null default '0' unique,
    url_token_code  int(5)                    default 0,
    create_at       datetime         not null default current_timestamp,
    following_count int(11) unsigned          default 0,
    followed_count  int(11) unsigned          default 0,
    marked_count    int(11) unsigned          default 0,
    question_count  int(11) unsigned          default 0,
    answer_count    int(11) unsigned          default 0,
    primary key (id)
) default charset = utf8;
#
drop table if exists questions;
create table questions
(
    id             int(11) unsigned not null auto_increment,
    user_id        int(11) unsigned not null,
    title          varchar(100)     not null default '',
    detail         longtext,
    create_at      datetime         not null default current_timestamp,
    modified_at    datetime         not null default current_timestamp,
    mark_count     int(11) unsigned not null default 0,
    follower_count int(11) unsigned not null default 0,
    answer_count   int(11) unsigned not null default 0,
    primary key (id)
) default charset = utf8;
#
drop table if exists answers;
alter table questions rename column mark_count to comment_count;
create table answers
(
    id            int(11) unsigned not null auto_increment,
    user_id       int(11) unsigned not null,
    question_id   int(11) unsigned not null,
    content       longtext,
    create_at     datetime         not null default current_timestamp,
    modified_at   datetime         not null default current_timestamp,
    marked_count  int(11) unsigned          default 0,
    comment_count int(11) unsigned          default 0,
    is_deleted    tinyint(1)                default 0,
    primary key (id),
    unique (user_id, question_id)
) default charset = utf8;
#
drop table if exists topics;
create table topics
(
    id   int(11) unsigned not null auto_increment,
    name varchar(50)      not null default '' unique,
    primary key (id)
) default charset = utf8;
#
drop table if exists question_topics;
create table question_topics
(
    question_id int(11) unsigned not null,
    topic_id    int(11) unsigned not null,
    primary key (question_id, topic_id)
) default charset = utf8;
#
drop table if exists question_comments;
create table question_comments
(
    id          int(11) unsigned not null auto_increment,
    user_id     int(11) unsigned not null,
    question_id int(11) unsigned not null,
    create_at   datetime         not null default current_timestamp,
    content     longtext,
    primary key (id)
) default charset = utf8;
#
drop table if exists answer_comments;
create table answer_comments
(
    id        int(11) unsigned not null auto_increment,
    user_id   int(11) unsigned not null,
    answer_id int(11) unsigned not null,
    create_at datetime         not null default current_timestamp,
    content   longtext,
    primary key (id)
) default charset = utf8;
#
drop table if exists question_followers;
create table question_followers
(
    question_id int(11) unsigned not null,
    member_id   int(11) unsigned not null,
    create_at   datetime         not null default current_timestamp,
    primary key (question_id, member_id)
) default charset = utf8;
#
alter table question_followers rename column member_id to user_id;
#
drop table if exists answer_followers;
create table answer_followers
(
    answer_id   int(11) unsigned not null,
    follower_id int(11) unsigned not null,
    create_at   datetime         not null default current_timestamp,
    primary key (answer_id, follower_id)
) default charset = utf8;
#
drop table if exists answer_voters;
create table answer_voters
(
    id        int(11) unsigned not null auto_increment,
    answer_id int(11) unsigned not null,
    voter_id  int(11) unsigned not null,
    primary key (id)
) default charset = utf8;
#
drop table if exists member_followers;
create table member_followers
(
    member_id   int(11) unsigned not null,
    follower_id int(11) unsigned not null,
    create_at   datetime         not null default current_timestamp,
    primary key (member_id, follower_id)
) default charset = utf8;
#
insert users
set created_at=now(),
    email='root@root.com',
    fullname='root',
    password='5508c15e9a9781f58869db41eb17f5c3',
    headline='a man to look for answer',
    url_token='root',
    avatar_url='/static/favicon.ico',
    question_count=1;
#
insert users
set created_at=now(),
    email='qiu@qiu.com',
    fullname='qiu',
    password='d70fcd4ea1b4bfa8d3301f7efe6ac779',
    gender=1,
    headline='an excellent man',
    url_token='qiu',
    avatar_url='/static/images/20180428152318.jpg',
    follower_count=1;
#
insert users
set created_at=now(),
    email='pan@pan.com',
    fullname='pan',
    password='87c8e41d5a549565eba508f243f31d49',
    gender=1,
    headline='a handsome man',
    url_token='pan';
#
insert questions
set user_id=1,
    title='对于知乎新人你们有什么好的建议？',
    detail='我是一个知乎的新用户，各位大大们有什么好的关注话题推荐一下，或者一些好的建议，或者您认为有用的东西，都推荐一下。
我的好奇心比较重，所以关注的东西比较多，但我不介意再多一点',
    created_at=now(),
    modified_at=now(),
    follower_count=1,
    comment_count=1;
#
insert question_comments
set user_id=3,
    question_id=1,
    created_at=now(),
    content='good question!!';
#
insert topics
set name='知乎';
#
insert question_topics
set question_id=1,
    topic_id=1;
#
insert question_followers
set question_id=1,
    user_id=1;
#
insert member_followers
set member_id=2,
    follower_id=1;
#
insert into topics (name)
values ('足球'),
       ('篮球'),
       ('排球'),
       ('游戏'),
       ('优秀'),
       ('你好'),
       ('电脑'),
       ('编程');

