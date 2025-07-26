package domain

type RegisterUser struct {
	Login             string
	Email             string
	Phone             string
	Password          string
	PreferredCurrency CurrencyCode
	FirstName         string
	LastName          string
	Language          Language
	Timezone          string
}

func (m RegisterUser) ToUser() *User {
	return &User{
		Email:             m.Email,
		Login:             m.Login,
		Phone:             m.Phone,
		PreferredCurrency: m.PreferredCurrency,
		FirstName:         m.FirstName,
		LastName:          m.LastName,
		Language:          m.Language,
		Timezone:          m.Timezone,
	}
}

type LoginUser struct {
	Login    string // email, login, phone, public_id etc
	Password string
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
}
