import NotificationModel from '../models/notification';

const NotificationType = {
  NotiSlice: {
    system: () => {},
    replied: () => {},
    quoted: () => {},
    sliceInfo: () => {},
  },

  UnreadNotiCount: {
    system: () => {},
    replied: () => {},
    quoted: () => {},
  },

  SystemNoti: {
    id: () => {},
    type: () => {},
    eventTime: () => {},
    hasRead: () => {},
    title: () => {},
    content: () => {},
  },

  RepliedNoti: {
    id: () => {},
    type: () => {},
    eventTime: () => {},
    hasRead: () => {},
    thread: () => {},
    repliers: () => {},
  },

  QuotedNoti: {
    id: () => {},
    type: () => {},
    eventTime: () => {},
    hasRead: () => {},
    thread: () => {},
    quotedPost: () => {},
    post: () => {},
    quoter: () => {},
  },
};

const unreadNotiCount = (ctx) => {
  console.log(ctx);
  if (!ctx.user) return null;
  return ctx.user;
};

export default NotificationType;
export { unreadNotiCount };
