import ConfigModel from '~/models/config';

const Query = {
  config: (obj, args, ctx) => ctx.config,
};

const Mutation = {
  editConfig: async (obj, { config }) => {
    const newConfig = await ConfigModel.setConfig(config);
    return newConfig;
  },
};

// Auto Resolver:
// Config

export default {
  Query,
  Mutation,
};
