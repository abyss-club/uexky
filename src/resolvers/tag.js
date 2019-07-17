import TagModel from '~/models/tag';

const Query = {
  mainTags: () => TagModel.getMainTags(),
  recommended: () => TagModel.getMainTags(),
  tags: async (_obj, { query, limit = 10 }) => TagModel.findTags({ query, limit }),
};

// Default Types Resolvers:
//   Tag:
//     name, isMain, belongsTo

const Tag = {
  belongsTo: tag => tag.getBelongsTo(),
};

export default {
  Query,
  Tag,
};
