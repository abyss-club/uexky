// user aggragate: user

package entity

type UserRepo interface{}

type UserService struct {
	Repo UserRepo `wire:"-"` // TODO
}

func NewUserService(repo UserRepo) UserService {
	return UserService{repo}
}

type User struct {
	Email   string       `json:"email"`
	Name    *string      `json:"name"`
	Tags    []string     `json:"tags"`
	Role    *string      `json:"role"`
	Threads *ThreadSlice `json:"threads"`
	Posts   *PostSlice   `json:"posts"`
}
