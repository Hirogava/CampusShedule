package types

type LessonType string

const (
    Lecture  LessonType = "lecture"
    Seminar  LessonType = "seminar"
    Practice LessonType = "practice"
    Test     LessonType = "test"
    Exam     LessonType = "exam"
    Webinar  LessonType = "webinar"
)
