package models

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gomodule/redigo/redis"
	"log"
	"my-zhihu/utils"
	"strconv"
	"strings"
)

type Page struct {
	Paging
	sessions.Session
}
type Paging struct {
	IsStart bool   `json:"is_start"`
	IsEnd   bool   `json:"is_end"`
	Next    string `json:"next"`
}

func HomeTimeLine(uid uint) []*Action {
	if uid == 0 {
		return TopStory(uid)
	}
	key := fmt.Sprintf("home:%d", uid)
	conn := redisPool.Get()
	defer conn.Close()
	res, err := redis.Strings(conn.Do(zrevrange, key, 0, 9))
	if err != nil {
		log.Println("models.TopContent(): ", err)
		return nil
	}
	var actions []*Action
	for _, member := range res {
		action := new(Action)
		s := strings.SplitN(member, ":", 3)
		act, err := strconv.Atoi(s[1])
		id := s[2]
		if err != nil {
			continue
		}
		if id != "" {
			who, err := strconv.Atoi(s[0])
			if err != nil {
				continue
			}
			action.User = GetUserByID(who)
		}
		time, err := redis.Int64(conn.Do(zscore, key, member))
		if err != nil {
			log.Println("models.TopContent(): ", err)
			continue
		}
		action.DateCreated = utils.FormatBeforeUnixTime(time)
		switch act {
		case AskQuestionAction:
			action.Type = AskQuestionAction
			action.Question = GetQuestion(id, uid)
		case AnswerQuestionAction:
			action.Type = AnswerQuestionAction
			action.Answer = GetAnswer(id, uid, "before")
		case FollowQuestionAction:
			action.Type = FollowQuestionAction
			action.Question = GetQuestion(id, uid)
		case VoteUpAnswerAction:
			action.Type = VoteUpAnswerAction
			action.Answer = GetAnswer(id, uid, "before")
		default:
			action.Type = OtherAction
		}
		actions = append(actions, action)
	}
	//
	if len(res) == 0 {
		return TopStory(uid)
	}
	return actions
}

func TopStory(uid uint) []*Action {
	conn := redisPool.Get()
	defer conn.Close()
	var actions []*Action
	answers, err := redis.Strings(conn.Do(zrevrange, "rank", 0, 9))
	if err != nil {
		log.Println("models.TopContent(): ", err)
		return actions
	}
	for _, aid := range answers {
		action := new(Action)
		action.Type = OtherAction
		action.Answer = GetAnswer(aid, uid)
		actions = append(actions, action)
	}
	return actions
}

func GetQuestion(aid string, uid uint, options ...string) *Question {

}

func GetQuestionWithAnswers() {

}

func (p *Page) QuestionComments() {

}

func (p *Page) QuestionFollowers() {

}

func (p *Page) AnswerVoters(aid string, offset int, uid uint) []User {
	var voters = make([]User, 0)
	start := p.Session.Get("start" + aid)
	if start == nil {
		var newStart int
		if err := db.QueryRow("SELECT answer_voters.id FROM users, answer_voters WHERE users.id=answer_voters.user_id "+
			"AND answer_id=? ORDER BY answer_voters.id DESC LIMIT 1", aid).Scan(&newStart); err != nil {
			log.Println("*Page.AnswerVoters(): ", err)
			return voters
		}
		p.Session.Set("start"+aid, newStart)
		_ = p.Session.Save()
		start = newStart
		offset = 0
		p.Paging.IsStart = true
	}
	limit := fmt.Sprintf("limit %d,%d", offset, 10)
	rows, err := db.Query("SELECT users.id, users.fullname, users.gender, users.headline, "+
		"users.avatar_url, users.url_token, users.answer_count, users.follower_count FROM users, answer_voters "+
		"WHERE users.id=answer_voters.user_id AND answer_id=? AND answer_voters.id<=? ORDER BY answer_voters.id DESC "+limit,
		aid, start.(int))
	if err != nil {
		log.Println("*Page.AnswerVoters(): ", err)
		return voters
	}
	defer rows.Close()

	var i int
	for ; rows.Next(); i++ {
		var voter User
		if err := rows.Scan(&voter.ID, &voter.Name, &voter.Gender,
			&voter.Headline, &voter.AvatarURL, &voter.URLToken,
			&voter.AnswerCount, &voter.FollowerCount); err != nil {
			log.Println("*Page.AnswerVoters(): ", err)
			continue
		}
		voter.QueryRelationWithVisitor(uid)
		voters = append(voters, voter)
	}
	if i < 10 {
		p.Paging.IsEnd = true
	} else {
		p.Paging.Next = fmt.Sprintf("/api/answers/%s/voters?offset=%d", aid, offset+i)
	}
	return voters
}
