package models

import (
	"fmt"
	"log"
	"time"
)

type Action struct {
	*User       `json:"user"`
	Type        int `json:"type"`
	*Question   `json:"question"`
	*Answer     `json:"answer"`
	DateCreated string `json:"date_created"`
}

func HandleNewAction(uid uint, action int, id string) {
	var err error
	defer func() {
		if err != nil {
			log.Println("models.HandleNewAction err:", err)
		}
	}()
	conn := redisPool.Get()
	defer conn.Close()
	rows, err := db.Query("SELECT follower_id FROM member_followers where member_id = ?", uid)
	if err != nil {
		return
	}
	key := fmt.Sprintf("profile:%d", uid)
	field := fmt.Sprintf("%d:%s", action, id)
	_, _ = conn.Do(zadd, key, time.Now().Unix(), field)
	var fid string
	for rows.Next() {
		err = rows.Scan(&fid)
		if err != nil {
			continue
		}
		key := "home:" + fid
		now := time.Now().Unix()
		field := fmt.Sprintf("%d:%d:%s", uid, action, id)
		_, _ = conn.Do(zadd, key, now, field)
	}
	if action == VoteUpAnswerAction {
		n, err := conn.Do(scard, "upvoted:"+id)
		if err != nil {
			return
		}
		score := n.(int64)*432 + 86400
		println(id, score)
		_, err = conn.Do(zadd, "rank", score, id)
		if err != nil {
			return
		}
	}
}

func RemoveAction(uid uint, action int, id string) {
	var err error
	defer func() {
		if err != nil {
			log.Println("models.DeleteAction(): ", err)
		}
	}()

	conn := redisPool.Get()
	rows, err := db.Query("SELECT follower_id FROM member_followers WHERE member_id=?", uid)
	if err != nil {
		return
	}
	key := fmt.Sprintf("profile:%d", uid)
	member := fmt.Sprintf("%d:%s", action, id)
	_, _ = conn.Do(zrem, key, member)
	var fid string
	for rows.Next() {
		if err = rows.Scan(&fid); err != nil {
			continue
		}
		key := "home:" + fid
		member := fmt.Sprintf("%d:%d:%s", uid, action, id)
		_, _ = conn.Do(zrem, key, member)
	}
}
