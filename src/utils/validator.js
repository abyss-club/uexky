import Joi from '@hapi/joi';
import { ParamsError } from '~/utils/error';

function unicodeLength(str) {
  let len = 0;
  for (let i = 0; i < str.length; i += 1) {
    const w = (str.codePointAt(i) > 255) ? 2 : 1;
    len += w;
  }
  return len;
}

function isUnicodeLength(str, { min = 0, max }) {
  if (Object.prototype.toString.call(str) !== '[object String]') throw new ParamsError('Invalid string.');
  if (!Number.isInteger(min)) throw new ParamsError('Invalid min value.');
  if (!Number.isInteger(max) && typeof max !== 'undefined') throw new ParamsError('Invalid max value.');

  const len = unicodeLength(str);
  if (!max) return len >= min;
  return len >= min && len <= max;
}

const emailSchema = Joi.object().keys({
  email: Joi.string().email().required(),
});

function isValidEmail(email) {
  const { error } = emailSchema.validate({ email });
  if (error) {
    return false;
  }
  return true;
}

export default {
  isUnicodeLength,
  isValidEmail,
};
