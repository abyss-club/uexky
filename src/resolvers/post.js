import PostModel from '~/models/post';

const Query = {
  post: async (obj, { id }, ctx) => {
    await ctx.limiter.take(1);
    const post = await PostModel.findById(id).exec();
    return post;
  },
};

const Mutation = {
  pubPost: async (obj, { post }, ctx) => {
    const limiterCfg = await ctx.config.getRateLimit();
    await ctx.limiter.take(limiterCfg.costSchema.pubPost);
    const newPost = await PostModel.pubPost(ctx, post);
    return newPost;
  },
};

const Post = {
  content: post => post.getContent(),
  quotes: async (post, args, ctx) => {
    await ctx.limiter.take(post.quoteSuids.length);
    const quotes = await post.getQuotes();
    return quotes;
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
