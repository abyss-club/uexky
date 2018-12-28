import UserModel from '../models/user';
import TagModel from '../models/tag';

const resolvers = {
  Tags: () => ({}),

  TagTreeNode: {
    mainTag: () => {},
    subTags: () => {},
  },

  Query: {
    // tags: (_, __, ctx) => ({
    //   if (!ctx)
    // }),
  },
};

export default resolvers;
