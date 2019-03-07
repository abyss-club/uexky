import mongoose from 'mongoose';

import Uid from '~/uid';

const SchemaObjectId = mongoose.ObjectId;

// MODEL: UserAid
//        used for save anonymousId for user in threads.
const UserAidSchema = new mongoose.Schema({
  userId: SchemaObjectId,
  threadSuid: String,
  anonymousId: String, // format: Uid
});

UserAidSchema.statics.getAid = async function getAid(userId, threadSuid) {
  const result = await UserAidModel.findOneAndUpdate({
    userId, threadSuid,
  }, {
    $setOnInsert: {
      userId,
      threadSuid,
      anonymousId: Uid.decode(await Uid.newSuid()),
    },
    $set: { updatedAt: Date() },
  }, { new: true, upsert: true }).exec();
  return result.anonymousId;
};

const UserAidModel = mongoose.model('UserAid', UserAidSchema);

export default UserAidModel;
