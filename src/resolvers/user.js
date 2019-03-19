import UserModel, { ensureSignIn } from '../models/user';
import AuthModel from '../models/auth';

const Query = {
  profile: (_, __, ctx) => ensureSignIn(ctx),
};

const Mutation = {
  auth: async (obj, { email }, ctx) => {
    if (ctx.user) throw new Error('Already signed in.');
    await AuthModel(ctx).addToAuth({ email });
    return true;
  },
  setName: async (obj, { name }, ctx) => {
    const user = ensureSignIn(ctx);
    const result = await UserModel(ctx).methods(user).setName(name);
    return result;
  },
  syncTags: async (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    const result = await UserModel(ctx).methods(user).syncTags(tags);
    return result;
  },
  addSubbedTags: async (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    const result = await UserModel(ctx).methods(user).addSubbedTags(tags);
    return result;
  },
  delSubbedTags: async (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    const result = await UserModel(ctx).methods(user).delSubbedTags(tags);
    return result;
  },

  // admin's apis:
  banUser: async (obj, { postId }, ctx) => {
    const user = ensureSignIn(ctx);
    await UserModel(ctx).methods(user).banUser(postId);
  },
  blockPost: async (obj, { postId }, ctx) => {
    const user = ensureSignIn(ctx);
    await UserModel(ctx).methods(user).blockPost(postId);
  },
  lockThread: async (obj, { threadId }, ctx) => {
    const user = ensureSignIn(ctx);
    await UserModel(ctx).methods(user).lockThread(threadId);
  },
  blockThread: async (obj, { threadId }, ctx) => {
    const user = ensureSignIn(ctx);
    await UserModel(ctx).methods(user).blockThread(threadId);
  },
  editTags: async (obj, { threadId, mainTag, subTags }, ctx) => {
    const user = ensureSignIn(ctx);
    await UserModel(ctx).methods(user).editTags(threadId, mainTag, subTags);
  },
};

// Default Types Resolver:
//   User:
//     email, name, tags

export default {
  Query,
  Mutation,
};
