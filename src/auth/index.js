import env from '~/utils/env';
import log from '~/utils/log';

import Code, { expireTime } from './code';
import Token from './token';

function authMiddleware(endpoint) {
  return async (ctx, next) => {
    const token = ctx.cookies.get('token') || '';
    if ((ctx.url === endpoint) && (token !== '')) {
      try {
        const email = await Token.getEmailByToken(token);
        const user = await UserModel.getUserByEmail(email, true);
        ctx.user = user;
        setCookie(ctx, token);
      } catch (e) {
        if (e.authError) ctx.user = null;
        else throw new Error(e);
      }
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
        setCookie(ctx, token);
        ctx.response.header.set('Location', `${env.PROTO}://${env.DOMAIN}`);
        ctx.response.header.set('Cache-Control', 'no-cache, no-store');
        ctx.response.status = 302;
      } catch (e) {
        log.error(e);
        ctx.throw(401, '验证信息错误或已失效');
      }
    }
    await next();
  };
}

function setCookie(ctx, token) {
  const opts = {
    path: '/',
    domain: env.DOMAIN,
    maxAge: expireTime.token,
    httpOnly: true,
    overwrite: true,
  };
  if (env.PROTO === 'https') {
    opts.secure = true;
  }
  ctx.cookies.set('token', token, opts);
}


export { authMiddleware, authHandler };
