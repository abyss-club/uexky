import ConfigModel from '~/models/config';

const Query = {
  config: (_obj, _args, ctx) => ctx.config,
};

const Mutation = {
  editConfig: async (_obj, { config }, ctx) => {
    const newConfig = await ConfigModel.setConfig(ctx, config);
    return newConfig;
  },
};

// Auto Resolver:
// Config

export default {
  Query,
  Mutation,
};
