import Mailgun from 'mailgun-js';

import log from '~/utils/log';
import env from '~/utils/env';

const mailgun = new Mailgun({
  apiKey: env.MAILGUN_PRIVATE_KEY, domain: env.MAILGUN_DOMAIN,
});

export default ({ sendTo, authCode }) => {
  const mailData = {
    // Specify email data
    from: env.MAILGUN_SENDER,
    // The email to contact
    to: sendTo,
    // Subject and text data
    subject: 'Hello from Mailgun',
    html: `Hello, This is not a plain-text email, I wanted to test some spicy Mailgun sauce in NodeJS! Authcode is ${authCode} .`,
  };
  mailgun.messages().send(mailData, (err, body) => {
    // If there is an error, render the error page
    if (err) {
      log.error(err);
    } else {
      // Here "submitted.jade" is the view file for this landing page
      // We pass the variable "email" from the url parameter in an object rendered by Jade
      log.info('send message', body);
    }
  });
};
