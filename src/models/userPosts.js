import mongoose from 'mongoose';

const SchemaObjectId = mongoose.ObjectId;

// MODEL: UserPosts
//        used for querying user's posts grouped by thread.
const UserPostsSchema = new mongoose.Schema({
  userId: SchemaObjectId,
  threadSuid: String,
  updatedAt: Date,
  posts: [{
    suid: String,
    createdAt: Date,
    anonymous: Boolean,
  }],
}, { autoCreate: true });

const UserPostsModel = mongoose.model('UserPosts', UserPostsSchema);

export default UserPostsModel;
