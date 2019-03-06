import ThreadModel from '~/models/thread';

const Query = {
  threadSlice: async (obj, { tags, query }, ctx) => {
    await ctx.limiter.take(query.limit);
    const threadSlice = await ThreadModel.getThreadSlice(tags, query);
    return threadSlice;
  },

  thread: async (obj, { id }, ctx) => {
    await ctx.limiter.take(1);
    const thread = await ThreadModel.findByUid(id);
    return thread;
  },
};

const Mutation = {
  pubThread: async (obj, { thread }, ctx) => {
    const { rateCost } = ctx.config;
    await ctx.limiter.take(rateCost.pubThread);
    const newThread = await ThreadModel.pubThread(ctx, thread);
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
  id: thread => thread.uid(),
  content: thread => thread.getContent(),
  replies: async (thread, query, ctx) => {
    await ctx.limiter.take(query.limit);
    const result = await thread.replies(query);
    return result;
  },
};

export default {
  Query,
  Mutation,
  Thread,
};
