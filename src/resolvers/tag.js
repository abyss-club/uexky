import TagModel from '~/models/tag';

const Query = {
  tags: async () => {
    const mainTags = await TagModel().getMainTags();
    return { mainTags };
  },
};

// Default Types Resolvers:
//   TagTreeNode:
//     mainTag, subTags

const Tags = {
  mainTags: ({ mainTags }) => mainTags,
  tree: async (obj, { query, limit = 10 }, ctx) => {
    await ctx.limiter.take(limit);
    try {
      const tree = await TagModel(ctx).getTree(limit, query);
      return tree;
    } catch (e) {
      throw e;
    }
  },
};

export default {
  Query,
  Tags,
};
