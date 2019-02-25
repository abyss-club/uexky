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
  mainTags: async (obj, args, ctx) => {
    const config = await ctx.config.getMainTags();
    return config;
  },
  rateLimit: async (obj, args, ctx) => {
    const config = await ctx.config.getRateLimit();
    return JSON.stringify(config);
  },
};

export default {
  Query,
  Mutation,
  Config,
};
