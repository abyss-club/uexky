import ConfigModel from '~/models/config';

const Query = {
  config: (obj, args, ctx) => ctx.config,
};

const Mutation = {
  editConfig: async (obj, { config }) => {
    await ConfigModel.setConfig(config);
    return {};
  },
};

// Auto Resolver:
// Config

export default {
  Query,
  Mutation,
};
