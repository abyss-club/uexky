import ConfigModel from '~/models/config';

const Query = {
  config: (obj, args, ctx) => {
    console.log({ config: ctx.config });
    return ctx.config;
  },
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
