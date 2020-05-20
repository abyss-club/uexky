package server

import "net/http"

/*
function authMiddleware(endpoint) {
  return async (ctx, next) => {
    const token = ctx.cookies.get('token') || '';
    let email;
    if ((ctx.url === endpoint) && (token !== '')) {
      email = await Token.getEmailByToken(token, true);
    }
    if (email) {
      ctx.auth = await UserModel.authContext({ email });
      ctx.response.set({ 'Set-Cookie': genCookie(token) });
    } else {
      ctx.auth = await UserModel.authContext({});
    }
    await next();
  };
}
*/

// func (s *Server) withLog(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	})
// }

func (s *Server) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil && err != http.ErrNoCookie {
			writeError(w, err)
			return
		}
		var tok string
		if tokenCookie != nil {
			tok = tokenCookie.Value
		}
		ctx, err := s.service().CtxWithUserByToken(r.Context(), tok)
		if err != nil {
			writeError(w, err)
			return
		}
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
