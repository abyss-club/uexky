import TagModel from '~/models/tag';

const Query = {
  mainTags: () => TagModel.getMainTags(),
  tag: (_obj, { query, limit = 10 }) => TagModel.findTags({ query, limit }),
};

// Default Types Resolvers:
//   Tag:
//     name, isMain, belongsTo

export default {
  Query,
};
