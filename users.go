package main

type Users struct {
	Id             uint32 `gorm:"column:id;type:INT(11) UNSIGNED;PRIMARY_KEY;AUTO_INCREMENT;NOT NULL"`
	UserId         uint32 `gorm:"column:user_id;type:INT(10) UNSIGNED;NOT NULL"`
	Gender         uint8  `gorm:"column:gender;type:TINYINT(2) UNSIGNED;NOT NULL"`
	Avatar         string `gorm:"column:avatar;type:VARCHAR(50);NOT NULL"`
	FollowingCount int32  `gorm:"column:following_count;type:INT(10);NOT NULL"`
	FollowedCount  int32  `gorm:"column:followed_count;type:INT(10);NOT NULL"`
	MarkedCount    uint32 `gorm:"column:marked_count;type:INT(10) UNSIGNED;NOT NULL"`
	Email          string `gorm:"column:email;type:VARCHAR(20);NOT NULL"`
	Fullname       string `gorm:"column:fullname;type:VARCHAR(20);NOT NULL"`
	Password       string `gorm:"column:password;type:VARCHAR(20);NOT NULL"`
	Headline       string `gorm:"column:headline;type:VARCHAR(10);NOT NULL"`
}
