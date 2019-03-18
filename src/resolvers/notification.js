import NotificationModel, { notiTypes } from '~/models/notification';
import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';

const Query = {
  unreadNotiCount: (obj, args, ctx) => {
    if (!ctx.user) return null;
    return {};
  },
  notification: async (obj, { type, query }, ctx) => {
    if (!ctx.user) return null;
    await ctx.limiter.take(query.limit);
    const result = await NotificationModel().getNotiSlice(ctx.user, type, query);
    return result;
  },
};

const UnreadNotiCount = (function makeUnreadNotiCount() {
  const resolver = {};
  notiTypes.forEach((type) => {
    resolver[type] = async (obj, args, ctx) => {
      const { user, limiter } = ctx;
      await limiter.take(1);
      const count = await NotificationModel().getUnreadCount(user, type);
      return count;
    };
  });
  return resolver;
}());

// Default Types resolvers:
// SystemNoti, RepliedNoti, QuotedNoti:
//     id, type, eventTime
const baseNoti = {
  hasRead: (noti, args, ctx) => {
    const { user } = ctx;
    return noti.eventTime > user.readNotiTime;
  },
};

const SystemNoti = Object.assign({
  title: noti => noti.system.title,
  content: noti => noti.system.content,
}, baseNoti);

const RepliedNoti = Object.assign({
  thread: async (noti, args, ctx) => {
    await ctx.limiter.take(1);
    const thread = await ThreadModel(ctx)
      .findByUid(noti.replied.threadId);
    return thread;
  },
  repliers: noti => noti.replied.repliers,
}, baseNoti);

const QuotedNoti = Object.assign({
  thread: async (noti, args, ctx) => {
    await ctx.limiter.take(1);
    const thread = await ThreadModel(ctx).findByUid(noti.quoted.threadId);
    return thread;
  },
  quotedPost: async (noti, args, ctx) => {
    await ctx.limiter.take(1);
    const post = await PostModel(ctx).findByUid(noti.quoted.quotedPostId);
    return post;
  },
  post: async (noti, args, ctx) => {
    await ctx.limiter.take(1);
    const post = await PostModel(ctx).findByUid(noti.quoted.quotedPostId);
    return post;
  },
  quoter: noti => noti.quoted.quoter,
}, baseNoti);

export default {
  Query,
  UnreadNotiCount,
  SystemNoti,
  RepliedNoti,
  QuotedNoti,
};
