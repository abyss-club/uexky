package auth

import (
	"context"

	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type RepoImpl struct{}

func (r *RepoImpl) SetCode(ctx context.Context, email string, code Code) error {
	panic("not implemented") // TODO: Implement
}

func (r *RepoImpl) GetCodeEmail(ctx context.Context, code Code) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (r *RepoImpl) DelCode(ctx context.Context, code Code) error {
	panic("not implemented") // TODO: Implement
}

func (r *RepoImpl) GetUserByAuthInfo(ctx context.Context, ai AuthInfo) (*entity.User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *RepoImpl) GetToken(ctx context.Context, tok string) (*Token, error) {
	panic("not implemented") // TODO: Implement
}

func (r *RepoImpl) SetToken(ctx context.Context, token *Token) error {
	panic("not implemented") // TODO: Implement
}

// func (u *UserRepo) SetCode(ctx context.Context, email string, code entity.Code) error {
// 	_, err := u.Redis.Set(string(code), email, entity.CodeExpire).Result()
// 	return redisErrWrapf(err, "SetCode(email=%s, code=%s)", email, code)
// }
//
// func (u *UserRepo) GetCodeEmail(ctx context.Context, code entity.Code) (string, error) {
// 	email, err := u.Redis.Get(string(code)).Result()
// 	return email, redisErrWrapf(err, "GetCodeEmail(code=%s)", code)
// }
//
// func (u *UserRepo) DelCode(ctx context.Context, code entity.Code) error {
// 	_, err := u.Redis.Del(string(code)).Result()
// 	return redisErrWrapf(err, "DelCode(code=%s)", code)
// }

// func (u *UserRepo) GetUserByAuthInfo(ctx context.Context, ai entity.AuthInfo) (*entity.User, error) {
// 	var user User
// 	if ai.IsGuest {
// 		if ai.UserID == 0 {
// 			return nil, uerr.New(uerr.ParamsError, "cannot get guest user without id")
// 		}
// 		data, err := u.Redis.Get(u.userRedisKey(ai.UserID)).Result()
// 		if err != nil {
// 			return nil, redisErrWrapf(err, "GetUserByAuthInfo(ai=%+v)", ai)
// 		}
// 		if err := json.Unmarshal([]byte(data), &user); err != nil {
// 			return nil, uerr.Wrapf(uerr.InternalError, err, "unmarshal user: %s", data)
// 		}
// 	} else {
// 		q := db(ctx).Model(&user)
// 		switch {
// 		case ai.UserID != 0:
// 			q = q.Where("id = ?", ai.UserID)
// 		case ai.Email != "":
// 			q = q.Where("email = ?", ai.Email)
// 		default:
// 			return nil, uerr.New(uerr.ParamsError, "cannot get signed user without id and email")
// 		}
// 		if err := q.Select(); err != nil {
// 			return nil, dbErrWrapf(err, "GetOrInsertUser.GetUser(ai=%+v)", ai)
// 		}
// 	}
// 	return u.toEntityUser(&user), nil
//
// }

// func (u *UserRepo) SetToken(ctx context.Context, token *entity.Token) error {
// 	data, err := json.Marshal(token)
// 	if err != nil {
// 		return uerr.Errorf(uerr.PermissionError, "SetToken(token=%+v), marshal json", token)
// 	}
// 	_, err = u.Redis.Set(token.Tok, data, token.Expire).Result()
// 	return redisErrWrapf(err, "SetToken(token=%+v)", token)
// }
//
// func (u *UserRepo) GetToken(ctx context.Context, tok string) (*entity.Token, error) {
// 	data, err := u.Redis.Get(tok).Result()
// 	if err != nil {
// 		return nil, redisErrWrapf(err, "GetToken(tok=%s)", tok)
// 	}
// 	var token entity.Token
// 	if err := json.Unmarshal([]byte(data), &token); err != nil {
// 		return nil, uerr.Wrapf(uerr.InternalError, err, "GetToken(tok=%s) unmarshal json: %s", tok, data)
// 	}
// 	return &token, nil
// }
