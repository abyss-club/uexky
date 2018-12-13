import mongoose from 'mongoose';
import { encode } from './uuid';

const UserSchema = new mongoose.Schema({
  email: { type: String, unique: true },
  name: { type: String, unique: true },
  tags: [String],
  read_noti_time: {
    system: Date,
    replied: Date,
    quoted: Date,
  },
  role: {
    type: String,
    range: [String],
  },
});
const UserModel = mongoose.model('User', UserSchema);
UserSchema.methods.anonymousId = async function anonymousId(threadId) {
  const obj = { userId: this.ObjectId, threadId };
  await UserAIDModel.update(obj, obj, { upsert: true });
  const aid = await UserAIDModel.findOne(obj);
  return aid.anonymousId;
};

const UserAIDSchema = new mongoose.Schema({
  userId: mongoose.Schema.Types.ObjectId,
  threadId: mongoose.Schema.Types.ObjectId,
});
const UserAIDModel = mongoose.model('UserAID', UserAIDSchema);
UserAIDSchema.methods.anonymousId = function anonymousId() {
  return encode(this.ObjectId);
};

export default UserModel;
export { UserAIDModel };
