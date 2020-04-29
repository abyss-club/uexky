package repo

import "gitlab.com/abyss.club/uexky/service"

func NewRepository() service.Repository {
	return service.Repository{
		User: &UserRepository{},
	}
}
