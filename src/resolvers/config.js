import ConfigModel from '~/models/config';

const Query = {
  config: () => ({}),
};

const Mutation = {
  editConfig: async (obj, { config }) => {
    if (config.mainTags) {
      await ConfigModel.modifyMainTags(config.mainTags);
    }
    if (config.rateLimit) {
      await ConfigModel.modifyRateLimit(config.rateLimit);
    }
    return {};
  },
};

const Config = {
  mainTags: (obj, args, ctx) => ctx.config.getMainTags(),
  rateLimit: (obj, args, ctx) => ctx.config.getRateLimit(),
};

export default {
  Query,
  Mutation,
  Config,
};
