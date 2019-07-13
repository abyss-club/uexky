import UserModel from '~/models/user';
import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';
import Code from '~/auth/code';
import { ParamsError } from '~/utils/error';

const Query = {
  profile: (_obj, _args, ctx) => ctx.auth.signedInUser(),
};


const Mutation = {
  auth: async (_obj, { email }, ctx) => {
    if (ctx.auth.isSignedIn) throw new Error('Already signed in.');
    await Code.addToAuth(email);
    return true;
  },
  setName: async (_obj, { name }, ctx) => {
    await UserModel.setName({ ctx, name });
    return ctx.auth.signedInUser();
  },
  syncTags: async (_obj, { tags }, ctx) => {
    await UserModel.syncTags({ ctx, tags });
    return ctx.auth.signedInUser();
  },
  addSubbedTag: async (_obj, { tag }, ctx) => {
    await UserModel.addSubbedTag({ ctx, tag });
    return ctx.auth.signedInUser();
  },
  delSubbedTag: async (_obj, { tag }, ctx) => {
    await UserModel.delSubbedTag({ ctx, tag });
    return ctx.auth.signedInUser();
  },

  // mod's apis:
  banUser: async (_obj, { postId, threadId }, ctx) => {
    let id;
    if (!postId) {
      const { userId } = await PostModel.findById({ postId });
      id = userId;
    } else if (!threadId) {
      const { userId } = await ThreadModel.findById({ threadId });
      id = userId;
    } else {
      throw new ParamsError('postId and threadId are both empty');
    }
    await UserModel.banUser({ ctx, userId: id });
    return true;
  },

  blockPost: async (_obj, { postId }, ctx) => {
    await PostModel.block({ ctx, postId });
    const post = await PostModel.findById({ postId });
    return post;
  },

  lockThread: async (_obj, { threadId }, ctx) => {
    await ThreadModel.lock({ ctx, threadId });
    const thread = await ThreadModel.findById({ threadId });
    return thread;
  },

  blockThread: async (_obj, { threadId }, ctx) => {
    await ThreadModel.block({ ctx, threadId });
    const thread = await ThreadModel.findById({ threadId });
    return thread;
  },

  editTags: async (_obj, { threadId, mainTag, subTags }, ctx) => {
    await ThreadModel.editTags({
      ctx, threadId, mainTag, subTags,
    });
    const thread = await ThreadModel.findById({ threadId });
    return thread;
  },
};

const User = {
  // auto field resolvers: email, name, role
  tags: user => user.getTags(),
  threads: (user, { query }) => ThreadModel.findUserThreads({ user, query }),
  posts: (user, { query }) => PostModel.findUserPosts({ user, query }),
};

export default {
  Query,
  Mutation,
  User,
};
