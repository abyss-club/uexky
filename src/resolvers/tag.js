import UserModel from '../models/user';
import TagModel from '../models/tag';
import config from '~/config';

const resolvers = {
  Tags: {
    mainTags: ({ mainTags }) => mainTags,
    tree: async (_, { _query, limit = 10 }) => {
      try {
        const tree = await TagModel.getTree(limit);
        return tree;
      } catch (e) {
        throw e;
      }
    },
  },

  Query: {
    tags: () => ({ mainTags: config.mainTags }),
  },
};

export default resolvers;
