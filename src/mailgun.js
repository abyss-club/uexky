import Mailgun from 'mailgun-js';

const mailData = {
  // Specify email data
  from: process.env.MAILGUN_SENDER,
  // The email to contact
  to: 'some@example.com',
  // Subject and text data
  subject: 'Hello from Mailgun',
  html: 'Hello, This is not a plain-text email, I wanted to test some spicy Mailgun sauce in NodeJS!',
};

const mailgun = new Mailgun({
  apiKey: process.env.MAILGUN_PRIVATE_KEY, domain: process.env.MAILGUN_DOMAIN,
});

export default () => {
  mailgun.messages().send(mailData, (err, body) => {
    // If there is an error, render the error page
    if (err) {
      console.log('got an error: ', err);
    } else {
      // Here "submitted.jade" is the view file for this landing page
      // We pass the variable "email" from the url parameter in an object rendered by Jade
      console.log(body);
    }
  });
};
