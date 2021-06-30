package repo

type CookiesRepository interface {
	InitTables()
	UpdateCookie(pt_pin, pt_key, usercookie string) (*Cookies, error)
	GetCookieByPtPin(pt_pin string) (*Cookies, error)
	DeleteCookieByPtPin(pt_pin string) (int64, error)
}
