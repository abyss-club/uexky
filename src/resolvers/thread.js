import ThreadModel from '../models/thread';

const ThreadTypes = {
  Thread: {
    id: () => {},
    anonymous: () => {},
    author: () => {},
    content: () => {},
    createdAt: () => {},
    mainTag: () => {},
    subTags: () => {},
    title: () => {},
    replies: () => {},
    replyCount: () => {},
    catalog: () => {},
  },

  ThreadCatalogItem: {
    postId: () => {},
    createdAt: () => {},
  },

  ThreadSlice: {
    threads: () => {},
    sliceInfo: () => {},
  },
};

const threadSlice = (ctx) => {
  console.log(ctx);
  if (!ctx.user) return null;
  return ctx.user;
};

const thread = (ctx) => {
  console.log(ctx);
  if (!ctx.user) return null;
  return ctx.user;
};

export default ThreadTypes;
export { threadSlice, thread };
