import { init } from '~/external/mailgun';

const mockMailgun = () => {
  init({
    mail: null,
    messages() {
      return (mail, fallback) => {
        this.mail = mail;
        fallback('success');
      };
    },
  });
};

export default mockMailgun;
