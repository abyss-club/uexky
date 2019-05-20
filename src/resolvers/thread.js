import ThreadModel from '~/models/thread';

const Query = {
  threadSlice: async (obj, { tags, query }, ctx) => {
    await ctx.limiter.take(query.limit);
    const threadSlice = await ThreadModel(ctx).getThreadSlice(tags, query);
    return threadSlice;
  },

  thread: async (obj, { id }, ctx) => {
    await ctx.limiter.take(1);
    const thread = await ThreadModel(ctx).findByUid(id);
    return thread;
  },
};

const Mutation = {
  pubThread: async (obj, { thread }, ctx) => {
    const { rateCost } = ctx.config;
    await ctx.limiter.take(rateCost.pubThread);
    const newThread = await ThreadModel(ctx).pubThread(thread);
    newThread.replyCount = 0;
    return newThread;
  },
};

// Default Types resolvers:
//   Thread:
//     id, anonymous, author, createdAt, mainTag,
//     subTags, title, replyCount, catelog
//   CatelogItem:
//     postId, createdAt,
//   ThreadSlice:
//     thread, sliceInfo
const Thread = {
  id: (thread, _, ctx) => ThreadModel(ctx).methods(thread).uid(),
  content: (thread, _, ctx) => ThreadModel(ctx).methods(thread).getContent(),
  replies: async (thread, query, ctx) => {
    await ctx.limiter.take(query.limit);
    const result = await ThreadModel(ctx).methods(thread).replies(query);
    return result;
  },
};

export default {
  Query,
  Mutation,
  Thread,
};
