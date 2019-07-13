import { mockMailgun } from '~/auth/mail';

export default () => {
  const mailgun = {
    mail: null,
    messages: function messages() {
      const that = this;
      return {
        send(mail, fallback) {
          that.mail = mail;
          fallback(null, 'success');
        },
      };
    },
  };
  mockMailgun(mailgun);
  return mailgun;
};
