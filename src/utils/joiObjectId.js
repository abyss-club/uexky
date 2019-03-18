import { ObjectId } from 'bson-ext';

module.exports = {
  name: 'objectId',
  pre(value, state, options) {
    if (!ObjectId.isValid(value)) {
      return this.createError('objectId.invalid', { value }, state, options);
    }

    if (options.convert) {
      return new ObjectId(value);
    }

    return value;
  },
  language: {
    invalid: 'Object needs to be a valid ObjectId.',
  },
};
