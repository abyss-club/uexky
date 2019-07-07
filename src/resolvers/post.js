import PostModel from '~/models/post';

const Query = {
  post: async (_obj, { id }, ctx) => {
    await ctx.limiter.take(1);
    const post = await PostModel.findById(id);
    return post;
  },
};

const Mutation = {
  pubPost: async (_obj, { post }, ctx) => {
    const { rateCost } = ctx.config;
    await ctx.limiter.take(rateCost.pubPost);
    const newPost = await PostModel.new({ ctx, post });
    return newPost;
  },
};

const Post = {
  // auto field resolvers: createdAt, anonymous, author, content, blocked
  id: post => post.id.duid,
  quotes: post => post.getQuotes(),
  quotedCount: post => post.getQuotedCount(),
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
