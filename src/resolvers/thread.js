import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';

const Query = {
  threadSlice: async (_obj, { tags, query }, ctx) => {
    await ctx.limiter.take(query.limit);
    const threadSlice = await ThreadModel.findSlice({ ctx, tags, query });
    return threadSlice;
  },

  thread: async (_obj, { id }, ctx) => {
    await ctx.limiter.take(1);
    const thread = await ThreadModel.findById({ id });
    return thread;
  },
};

const Mutation = {
  pubThread: async (_obj, { thread }, ctx) => {
    const { rateCost } = ctx.config;
    await ctx.limiter.take(rateCost.pubThread);
    const newThread = await ThreadModel.new({ ctx, thread });
    return newThread;
  },
};

const Thread = {
  // auto field resolvers: createdAt, anonymous, author, title, content, blocked, locked
  id: thread => thread.id.duid,
  mainTag: async thread => thread.getMainTag(),
  subTags: async thread => thread.getSubTags(),
  title: thread => thread.title,
  replies: async (thread, { query }, ctx) => {
    await ctx.limiter.take(query.limit);
    const result = await PostModel.findThreadPosts({ threadId: thread.id, query });
    return result;
  },
  replyCount: async (thread) => {
    const count = await PostModel.getThreadReplyCount({ threadId: thread.id });
    return count;
  },
  catalog: async (thread) => {
    const catalog = await PostModel.getThreadCatelog({ threadId: thread.id });
    return catalog;
  },
};

// auto type resolvers:
//   ThreadSlice
//   ThreadCatalogItem, with fields: postId, createdAt

export default {
  Query,
  Mutation,
  Thread,
};
