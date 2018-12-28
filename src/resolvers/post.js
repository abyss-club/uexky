import PostModel from '../models/post';

const PostTypes {
  Post: {
    id: () => {},
    anonymous: () => {},
    author: () => {},
    content: () => {},
    createTime: () => {},
    quotes: () => {},
    quoteCount: () => {},
  },

  PostSlice: {
    posts: () => {},
    sliceInfo: () => {},
  }
}

const post = (ctx) => {
  console.log(ctx);
  if (!ctx.user) return null;
  return ctx.user;
};

export default PostTypes;
export { post };
