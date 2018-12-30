import NotificationModel from '~/models/notification';
import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';

const Query = {
  unreadNotiCount: (obj, args, ctx) => {
    if (!ctx.user) return null;
    return ctx.user;
  },
  notification: async (obj, { type, query }, ctx) => {
    if (!ctx.user) return null;
    const result = await NotificationModel.getNotiSlice(ctx.user, type, query);
    return result;
  },
};

const getUnreadNotiCount = async (user, type) => {
  const count = await user.getUnreadNotiCount(type); // TODO
  return count;
};

const UnreadNotiCount = {
  system: user => getUnreadNotiCount(user, 'system'),
  replied: user => getUnreadNotiCount(user, 'replied'),
  quoted: user => getUnreadNotiCount(user, 'quoted'),
};

// Default Resolve Fields:
// id, type, eventTime
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
  thread: async (noti) => {
    const thread = await ThreadModel
      .findOne({ _id: noti.replied.threadId }).exec();
    return thread;
  },
  repliers: noti => noti.repliers,
}, baseNoti);

const QuotedNoti = Object.assign({
  thread: async (noti) => {
    const thread = await ThreadModel.findById(noti.replied.threadId).exec();
    return thread;
  },
  quotedPost: async (noti) => {
    const post = await PostModel.findById(noti.replied.quotedPostId).exec();
    return post;
  },
  post: async (noti) => {
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
