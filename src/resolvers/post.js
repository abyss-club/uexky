import PostModel from '~/models/post';

const Query = {
  post: async (obj, { id }) => {
    const post = await PostModel.findById(id).exec();
    return post;
  },
};

const Mutation = {
  pubPost: async (obj, { post }, ctx) => {
    const newPost = await PostModel.pubPost(ctx, post);
    return newPost;
  },
};

const Post = {
  content: post => post.getContent(),
  quotes: post => post.getQuotes(),
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
