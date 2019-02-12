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
    const result = await NotificationModel.getNotiSlice(ctx.user, type, query);
    return result;
  },
};

const UnreadNotiCount = (function makeUnreadNotiCount() {
  const resolver = {};
  notiTypes.forEach((type) => {
    resolver[type] = async (obj, args, ctx) => {
      const { user, limiter } = ctx;
      await limiter.take(1);
      const count = await NotificationModel.getUnreadCount(user, type);
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
    const thread = await ThreadModel
      .findOne({ _id: noti.replied.threadId }).exec();
    return thread;
  },
  repliers: noti => noti.repliers,
}, baseNoti);

const QuotedNoti = Object.assign({
  thread: async (noti, args, ctx) => {
    await ctx.limiter.take(1);
    const thread = await ThreadModel.findById(noti.replied.threadId).exec();
    return thread;
  },
  quotedPost: async (noti, args, ctx) => {
    await ctx.limiter.take(1);
    const post = await PostModel.findById(noti.replied.quotedPostId).exec();
    return post;
  },
  post: async (noti, args, ctx) => {
    await ctx.limiter.take(1);
    const post = await PostModel.findById(noti.replied.quotedPostId).exec();
    return post;
  },
  quoter: noti => noti.replied.quoter,
}, baseNoti);

export default {
  Query,
  UnreadNotiCount,
  SystemNoti,
  RepliedNoti,
  QuotedNoti,
};
