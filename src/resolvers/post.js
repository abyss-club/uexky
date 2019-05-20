import PostModel from '~/models/post';

const Query = {
  post: async (obj, { id }, ctx) => {
    await ctx.limiter.take(1);
    const post = await PostModel(ctx).findByUid(id);
    return post;
  },
};

const Mutation = {
  pubPost: async (obj, { post }, ctx) => {
    const { rateCost } = ctx.config;
    await ctx.limiter.take(rateCost.pubPost);
    const newPost = await PostModel(ctx).pubPost(post);
    newPost.quoteCount = 0;
    return newPost;
  },
};

const Post = {
  id: (post, _, ctx) => PostModel(ctx).methods(post).uid(),
  content: (post, _, ctx) => PostModel(ctx).methods(post).getContent(),
  quotes: async (post, args, ctx) => {
    if (post.quoteSuids) {
      await ctx.limiter.take(post.quoteSuids.length);
      const quotes = await PostModel(ctx).methods(post).getQuotes();
      return quotes;
    }
    return [];
  },
  quoteCount: async (post, args, ctx) => {
    const quoteCount = await PostModel(ctx).methods(post).quoteCount();
    return quoteCount;
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
