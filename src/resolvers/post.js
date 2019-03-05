import PostModel from '~/models/post';

const Query = {
  post: async (obj, { id }, ctx) => {
    await ctx.limiter.take(1);
    const post = await PostModel.findByUid(id);
    return post;
  },
};

const Mutation = {
  pubPost: async (obj, { post }, ctx) => {
    const { rateCost } = ctx.config;
    await ctx.limiter.take(rateCost.pubPost);
    const newPost = await PostModel.pubPost(ctx, post);
    return newPost;
  },
};

const Post = {
  id: post => post.uid(),
  content: post => post.getContent(),
  quotes: async (post, args, ctx) => {
    if (post.quoteSuids) {
      await ctx.limiter.take(post.quoteSuids.length);
      const quotes = await post.getQuotes();
      return quotes;
    }
    return [];
  },
};

// Default Types Resolvers:
//   Post:
//     id， anonymous, author, createdAt, quotes, quoteCount
//   PostSlice:
//     posts, sliceInfo

export default {
  Query,
  Mutation,
  Post,
};
