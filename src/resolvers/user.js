import { ensureSignIn } from '../models/user';
import AuthModel from '../models/auth';

const Query = {
  profile: (_, __, ctx) => ensureSignIn(ctx),
};

const Mutation = {
  auth: (obj, { email }, ctx) => {
    if (ctx.user) throw new Error('Already signed in.');
    AuthModel.addToAuth(email);
    return true;
  },
  setName: async (obj, { name }, ctx) => {
    const user = ensureSignIn(ctx);
    const result = await user.setName(name);
    return result;
  },
  syncTags: async (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    const result = await user.syncTags(tags);
    return result;
  },
  addSubbedTags: async (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    const result = await user.addSubbedTags(tags);
    return result;
  },
  delSubbedTags: async (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    const result = await user.delSubbedTags(tags);
    return result;
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
