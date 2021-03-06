package main

import pluginapi "github.com/mattermost/mattermost-plugin-api"

type Store interface {
	StoreQuiz(q *Quiz) error
	GetQuiz(id string) (*Quiz, error)
	DeleteQuiz(id string) error

	AddAvailableQuiz(q *Quiz) error
	GetAvailableQuizes() []*Quiz

	GetGame(id string) (*Game, error)
	StoreGame(g *Game) error
	DeleteGame(id string) error

	StoreCourse(c *Course) error
	GetCourse(id string) (*Course, error)
	AddAvailableCourse(c *Course) error
	DeleteCourse(id string) error
}

const (
	KVQuizPrefix   = "quiz_"
	KVQuizList     = "quizList"
	KVGamePrefix   = "game_"
	KVCoursePrefix = "course_"
	KVCourseList   = "courseList"
)

type store struct {
	mm *pluginapi.Client
}

func NewStore(mm *pluginapi.Client) Store {
	return &store{
		mm: mm,
	}
}

func (s *store) GetGame(id string) (*Game, error) {
	var g *Game
	err := s.mm.KV.Get(getGameKey(id), &g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (s *store) StoreGame(g *Game) error {
	_, err := s.mm.KV.Set(getGameKey(g.RootPostID), g)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) DeleteGame(id string) error {
	err := s.mm.KV.Delete(getGameKey(id))
	if err != nil {
		return err
	}

	return nil
}

func (s *store) GetAvailableQuizes() []*Quiz {
	out := []*Quiz{}

	quizIDList := []string{}
	err := s.mm.KV.Get(KVQuizList, &quizIDList)
	if err != nil {
		s.mm.Log.Debug("Cannot get quiz list", "error", err)
		return out
	}

	for _, id := range quizIDList {
		q, err := s.GetQuiz(id)
		if err != nil {
			s.mm.Log.Debug("Error getting quiz", "id", id, "err", err)
			continue
		}

		if q == nil {
			s.mm.Log.Debug("Quiz not found", "id", id)
			continue
		}

		out = append(out, q)
	}

	return out
}

func (s *store) AddAvailableQuiz(q *Quiz) error {
	quizIDList := []string{}
	err := s.mm.KV.Get(KVQuizList, &quizIDList)
	if err != nil {
		return err
	}

	for _, id := range quizIDList {
		if id == q.ID {
			return nil
		}
	}

	quizIDList = append(quizIDList, q.ID)

	_, err = s.mm.KV.Set(KVQuizList, quizIDList)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) removeAvailableQuiz(id string) error {
	quizIDList := []string{}
	err := s.mm.KV.Get(KVQuizList, &quizIDList)
	if err != nil {
		return err
	}

	changed := false
	for i, qid := range quizIDList {
		if qid == id {
			quizIDList = append(quizIDList[0:i], quizIDList[i+1:]...)
			changed = true
			break
		}
	}

	if !changed {
		return nil
	}

	_, err = s.mm.KV.Set(KVQuizList, quizIDList)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) DeleteQuiz(id string) error {
	err := s.removeAvailableQuiz(id)
	if err != nil {
		return err
	}

	err = s.mm.KV.Delete(getQuizKey(id))
	if err != nil {
		return err
	}

	return nil
}

func (s *store) StoreQuiz(q *Quiz) error {
	_, err := s.mm.KV.Set(getQuizKey(q.ID), q)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) GetQuiz(id string) (*Quiz, error) {
	q := &Quiz{}
	err := s.mm.KV.Get(getQuizKey(id), q)
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (s *store) StoreCourse(c *Course) error {
	_, err := s.mm.KV.Set(getCourseKey(c.ID), c)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) GetCourse(id string) (*Course, error) {
	c := &Course{}
	err := s.mm.KV.Get(getCourseKey(id), c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *store) AddAvailableCourse(c *Course) error {
	courseIDList := []string{}
	err := s.mm.KV.Get(KVCourseList, &courseIDList)
	if err != nil {
		return err
	}

	for _, id := range courseIDList {
		if id == c.ID {
			return nil
		}
	}

	courseIDList = append(courseIDList, c.ID)

	_, err = s.mm.KV.Set(KVCourseList, courseIDList)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) DeleteCourse(id string) error {
	err := s.removeAvailableCourse(id)
	if err != nil {
		return err
	}

	err = s.mm.KV.Delete(getCourseKey(id))
	if err != nil {
		return err
	}

	return nil
}

func (s *store) removeAvailableCourse(id string) error {
	courseIDList := []string{}
	err := s.mm.KV.Get(KVCourseList, &courseIDList)
	if err != nil {
		return err
	}

	changed := false
	for i, qid := range courseIDList {
		if qid == id {
			courseIDList = append(courseIDList[0:i], courseIDList[i+1:]...)
			changed = true
			break
		}
	}

	if !changed {
		return nil
	}

	_, err = s.mm.KV.Set(KVCourseList, courseIDList)
	if err != nil {
		return err
	}
	return nil
}

func getQuizKey(id string) string {
	return KVQuizPrefix + id
}

func getGameKey(id string) string {
	return KVGamePrefix + id
}

func getCourseKey(id string) string {
	return KVCoursePrefix + id
}
