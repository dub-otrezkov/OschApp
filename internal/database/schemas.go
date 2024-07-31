package db

type Submission struct {
	TaskId    int
	Verdict   int
	SessionId int
}

type Session struct {
	Id     int
	UserId int
	Active bool
	ExamId int
}
