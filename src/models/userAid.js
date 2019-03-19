import JoiBase from 'joi';
import JoiObjectId from '~/utils/joiObjectId';
import { ParamsError } from '~/utils/error';
import mongo from '~/utils/mongo';
import log from '~/utils/log';

import Uid from '~/uid';

const Joi = JoiBase.extend(JoiObjectId);
const USERAID = 'userAid';
const col = () => mongo.collection(USERAID);

const userAidSchema = Joi.object().keys({
  userId: Joi.objectId().required(),
  threadSuid: Joi.string().alphanum().length(15),
  anonymousId: Joi.string().alphanum(),
});

// MODEL: UserAid
//        used for save anonymousId for user in threads.
// const UserAidSchema = new mongoose.Schema({
//   userId: SchemaObjectId,
//   threadSuid: String,
//   anonymousId: String, // format: Uid
// });

const UserAidModel = () => ({
  getAid: async function getAid(userId, threadSuid) {
    const { value, error } = userAidSchema.validate({ userId, threadSuid });
    if (error) {
      log.error(error);
      throw new ParamsError(`Thread validation failed, ${error}`);
    }
    const result = await col().findOneAndUpdate(value, {
      $setOnInsert: {
        ...value,
        anonymousId: Uid.decode(await Uid.newSuid()),
      },
      $set: { updatedAt: Date() },
    }, { returnOriginal: false, upsert: true });
    return result.value.anonymousId;
  },
});

export default UserAidModel;
