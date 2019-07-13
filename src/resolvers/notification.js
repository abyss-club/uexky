import NotificationModel, { NOTI_TYPES } from '~/models/notification';
import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';

const Query = {
  unreadNotiCount: (_obj, _args, ctx) => {
    if (!ctx.user) return { system: 0, replied: 0, quoted: 0 };
    return {};
  },
  notification: async (_obj, { type, query }, ctx) => {
    await ctx.limiter.take(query.limit);
    const noti = await NotificationModel.findNotiSlice({ ctx, type, query });
    return noti;
  },
};

const UnreadNotiCount = {
  system: async (_ojb, _args, ctx) => {
    await ctx.limiter.take(1);
    const notis = await NotificationModel.getUnreadCount({ ctx, type: NOTI_TYPES.SYSTEM });
    return notis;
  },
  replied: async (_ojb, _args, ctx) => {
    await ctx.limiter.take(1);
    const notis = await NotificationModel.getUnreadCount({ ctx, type: NOTI_TYPES.REPLIED });
    return notis;
  },
  quoted: async (_ojb, _args, ctx) => {
    await ctx.limiter.take(1);
    const notis = await NotificationModel.getUnreadCount({ ctx, type: NOTI_TYPES.QUOTED });
    return notis;
  },
};

// Default Type resolvers:
//   SystemNoti: type, eventTime, hasRead, title, content
const SystemNoti = {
  id: noti => noti.key,
};

// Default Type resolvers:
//   RepliedNoti: type, eventTime, hasRead
const RepliedNoti = {
  id: noti => noti.key,
  thread: noti => ThreadModel.findById({ threadId: noti.threadId }),
  repliers: async (noti) => {
    const { posts } = await PostModel.findThreadPosts({
      threadId: noti.threadId,
      query: { before: '', limit: 5 },
    });
    return posts.map(p => p.author);
  },
};

// Default Type resolvers:
//   QuotedNoti: type, eventTime, hasRead, quoter
const QuotedNoti = {
  id: noti => noti.key,
  thread: noti => ThreadModel.findById({ threadId: noti.threadId }),
  quotedPost: noti => PostModel.findById({ postId: noti.quotedId }),
  post: noti => PostModel.findById({ postId: noti.postId }),
};

export default {
  Query,
  UnreadNotiCount,
  SystemNoti,
  RepliedNoti,
  QuotedNoti,
};
