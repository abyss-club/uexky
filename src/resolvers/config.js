import ConfigModel from '~/models/config';

const Query = {
  config: (obj, args, ctx) => ctx.config,
};

const Mutation = {
  editConfig: async (obj, { config }, ctx) => {
    const newConfig = await ConfigModel(ctx).setConfig(config);
    return newConfig;
  },
};

// Auto Resolver:
// Config

export default {
  Query,
  Mutation,
};
