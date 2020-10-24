package entity

type Repo struct {
	User   UserRepo
	Thread ThreadRepo
	Post   PostRepo
	Tag    TagRepo
	Noti   NotiRepo
}
