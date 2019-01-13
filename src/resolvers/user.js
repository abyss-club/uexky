import { ensureSignIn } from '../models/user';
import AuthModel from '../models/auth';

const Query = {
  profile: (_, __, ctx) => ensureSignIn(ctx),
};

const Mutation = {
  auth: (obj, { email }, ctx) => {
    if (ctx.user) throw new Error('already signed in!');
    AuthModel.addToAuth(email);
    return true;
  },
  setName: (obj, { name }, ctx) => {
    const user = ensureSignIn(ctx);
    user.setName(name);
  },
  syncTags: (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    user.syncTags(tags);
  },
  addSubbedTags: (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    user.addSubbedTags(tags);
  },
  delSubbedTags: (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    user.delSubbedTags(tags);
  },

  // admin's apis:
  banUser: async (obj, { postId }, ctx) => {
    const user = ensureSignIn(ctx);
    await user.banUser(postId);
  },
  blockPost: async (obj, { postId }, ctx) => {
    const user = ensureSignIn(ctx);
    await user.blockPost(postId);
  },
  lockThread: async (obj, { threadId }, ctx) => {
    const user = ensureSignIn(ctx);
    await user.lockThread(threadId);
  },
  blockThread: async (obj, { threadId }, ctx) => {
    const user = ensureSignIn(ctx);
    await user.blockThread(threadId);
  },
  editTags: async (obj, { threadId, mainTag, subTags }, ctx) => {
    const user = ensureSignIn(ctx);
    await user.editTags(threadId, mainTag, subTags);
  },
};

// Default Types Resolver:
//   User:
//     email, name, tags

export default {
  Query,
  Mutation,
};
