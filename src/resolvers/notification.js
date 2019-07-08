import NotificationModel, { NOTI_TYPES } from '~/models/notification';
import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';

const Query = {
  unreadNotiCount: (_obj, _args, ctx) => {
    if (!ctx.user) return null;
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
//   SystemNoti: id, type, eventTime, hasRead, title, content

// Default Type resolvers:
//   RepliedNoti: id, type, eventTime, hasRead
const RepliedNoti = {
  thread: noti => ThreadModel.findById(noti.threadId),
  repliers: noti => PostModel.getLastRepliers(noti.threadId),
};

// Default Type resolvers:
//   QuotedNoti: id, type, eventTime, hasRead, quoter
const QuotedNoti = {
  thread: noti => ThreadModel.findById(noti.threadId),
  quotedPost: noti => PostModel.findById(noti.quotedId),
  post: noti => PostModel.findById(noti.postId),
};

export default {
  Query,
  UnreadNotiCount,
  RepliedNoti,
  QuotedNoti,
};
