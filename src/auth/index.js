import UserModel from '~/models/user';
import env from '~/utils/env';
import log from '~/utils/log';

import Code, { expireTime } from './code';
import Token from './token';

function authMiddleware(endpoint) {
  return async (ctx, next) => {
    const token = ctx.cookies.get('token') || '';
    if ((ctx.url === endpoint) && (token !== '')) {
      const email = await Token.getEmailByToken(token);
      ctx.auth = await UserModel.authContext({ email });
      ctx.response.set({ 'Set-Cookie': genCookie(token) });
    } else {
      ctx.auth = await UserModel.authContext();
    }
    await next();
  };
}

function authHandler() {
  return async (ctx, next) => {
    if (!ctx.query.code || ctx.query.code.length !== 36) {
      ctx.throw(400, '验证信息格式错误');
    } else {
      try {
        const email = await Code.getEmailByCode(ctx.query.code);
        const token = await Token.genNewToken(email);
        ctx.response.set({
          Location: `${env.PROTO}://${env.DOMAIN}`,
          'Cache-Control': 'no-cache, no-store',
          'Set-Cookie': genCookie(token),
        });
        ctx.response.status = 302;
      } catch (e) {
        log.error(e);
        ctx.throw(401, '验证信息错误或已失效');
      }
    }
    await next();
  };
}

function genCookie(token) {
  const cookie = [
    `token=${token}`,
    ';Path=/',
    `;Max-Age=${expireTime.token}`,
    `;Domain=${env.DOMAIN}`,
    ';HttpOnly',
  ];
  if (env.PROTO === 'https') {
    cookie.push(';Secure');
  }
  return cookie.join('');
}


export { authMiddleware, authHandler };
