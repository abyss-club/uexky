import ThreadModel from '~/models/thread';
import Uid from '~/uid';

const Query = {
  threadSlice: async (obj, { tags, query }) => {
    const threadSlice = await ThreadModel.getThreadSlice(tags, query);
    return threadSlice;
  },

  thread: async (ctx, { id }) => {
    const thread = await ThreadModel.findById(Uid.encode(id)).exec();
    return thread;
  },
};

const Mutation = {
  pubThread: async (obj, { thread }, ctx) => {
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
  content: thread => thread.getContent(),
  replies: async (thread, { query }) => {
    const result = await thread.replies(query);
    return result;
  },
};

export default {
  Query,
  Mutation,
  Thread,
};
