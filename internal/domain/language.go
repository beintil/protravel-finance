package domain

type Language string

func (l Language) String() string {
	return string(l)
}

const (
	English Language = "en"
	Russian Language = "ru"
)

var Languages = []Language{
	English,
	Russian,
}

var LanguagesMap = map[Language]bool{
	English: true,
	Russian: true,
}
