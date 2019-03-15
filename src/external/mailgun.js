import Mailgun from 'mailgun-js';

import env from '~/utils/env';
import log from '~/utils/log';

let mailgun = null;

const init = (obj) => {
  if (!obj) {
    mailgun = new Mailgun({
      apiKey: env.MAILGUN_PRIVATE_KEY,
      domain: env.MAILGUN_DOMAIN,
    });
  } else {
    mailgun = obj;
  }
};

const sendAuthMail = (code, email) => new Promise(((resolve, reject) => {
  const codeUrl = `${env.API_DOMAIN}/auth/?code=${code}`;
  const mail = {
    from: `Abyss <auth@${env.MAILGUN_DOMAIN}>`,
    to: email,
    subject: '点击登入 Abyss!',
    text: `点击此链接进入 Abyss：${codeUrl}`,
    html: `<html>
  <head>
      <meta charset="utf-8">
      <title>点击登入 Abyss!</title>
  </head>
  <body>
      <p>点击 <a href="${codeUrl}">此链接</a> 进入 Abyss</p>
  </body>
</html>`,
  };
  if (!mailgun) {
    init();
  }
  mailgun.messages().send(mail, (err, body) => {
    if (err) {
      reject(err);
    } else {
      log.info(`send auth email to ${email}, res = ${body}`);
      resolve(body);
    }
  });
}));

export default sendAuthMail;
export { init };
