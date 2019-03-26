import { Base64 } from '~/uid';
import AuthModel from '~/models/auth';
import { mockMailgun } from '~/utils/authMail';

import startRepl from '../__utils__/mongoServer';


jest.setTimeout(60000);

let replSet;
let mongoClient;
let ctx;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

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


describe('Testing auth', () => {
  const authCode = Base64.randomString(36);
  const mockEmail = 'test@example.com';
  it('add user', async () => {
    const model = AuthModel(ctx);
    mockMailgun(mailgun);
    await model.addToAuth(mockEmail);
    const result = await model.col().findOne({ email: mockEmail });
    expect(result.email).toEqual(mockEmail);
    expect(mailgun.mail.to).toEqual(mockEmail);
    expect(mailgun.mail.text).toMatch(result.authCode);
  });
  it('validate user authCode for only once', async () => {
    const model = AuthModel(ctx);
    const doc = await model.col().findOne({ email: mockEmail });
    const result = await model.getEmailByCode(doc.authCode);
    expect(result).toEqual(mockEmail);
    const deletedResult = await model.col().findOne({ authCode });
    expect(deletedResult).toBeNull();
  });
});
