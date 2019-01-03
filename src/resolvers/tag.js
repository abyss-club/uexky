import TagModel from '~/models/tag';
import ConfigModel from '~/models/config';

const Query = {
  tags: () => ({ mainTags: ConfigModel.getMainTags() }),
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
