package db

type University struct {
	ID int
	Name string
}

type Group struct {
	ID int
	Name string
}

type LessonType string

const (
	Lecture LessonType = "lecture"
	Seminar LessonType = "seminar"
	Practice LessonType = "practice"
	Test LessonType = "test"
	Exam LessonType = "exam"
	Webinar LessonType = "webinar"
)

func (t LessonType) String() string {
	switch t {
	case Lecture:
		return "Лекция"
	case Seminar:
		return "Семинар"
	case Practice:
		return "Практика"
	case Exam:
		return "Экзамен"
	case Test:
		return "Зачёт"
	case Webinar:
		return "Вебинар"
	default:
		return "Другое"
	}
}

func (t LessonType) TypeToEmoji() string {
	switch t {
	case Lecture:
		return "📖"
	case Seminar:
		return "💬"
	case Practice:
		return "🧪"
	case Exam:
		return "🧾"
	case Test:
		return "✅"
	case Webinar:
		return "💻"
	default:
		return "📚"
	}
}

type Lesson struct {
	Teacher string
	Room string
	StartTime string
	EndTime string
	Date string
	DateOfWeek string
	Type LessonType
}

type Day struct {
	Lessons []Lesson
	WeekDay string
}
