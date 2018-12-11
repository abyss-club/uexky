import mongoose from 'mongoose';
import { encode } from './uuid';

const { ObjectId } = mongoose.Schema.Types;

const UserSchema = new mongoose.Schema({
  email: { type: String, unique: true },
  name: { type: String, unique: true },
  tags: [String],
  readNotiTime: {
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
UserSchema.methods.anonymousId = function anonymousId(threadId) {
  return threadId;
};

const UserAIDSchema = new mongoose.Schema({
  userId: ObjectId,
  threadId: ObjectId,
});
const UserAIDModel = mongoose.model('UserAID', UserAIDSchema);
UserAIDSchema.methods.anonymousId = function anonymousId() {
  return encode(this.ObjectId);
};

export default UserModel;
export { UserAIDModel };
