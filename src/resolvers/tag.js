import TagModel from '../models/tag';
import config from '~/config';

const Query = {
  tags: () => ({ mainTags: config.mainTags }),
};

// Default Types Resolvers:
//   TagTreeNode:
//     mainTag, subTags

const Tags = {
  mainTags: ({ mainTags }) => mainTags,
  tree: async (obj, { query, limit = 10 }) => {
    try {
      const tree = await TagModel.getTree(limit, query);
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
