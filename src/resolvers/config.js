import ConfigModel from '~/models/config';

const Query = {
  config: () => ({}),
};

const Mutation = {
  editConfig: async (obj, { config }) => {
    const returnConfig = {};
    if (config.mainTags) {
      const newMainTags = await ConfigModel.modifyMainTags(config.mainTags);
      returnConfig.mainTags = newMainTags;
    }
    if (config.rateLimit) {
      const newRateLimit = await ConfigModel.modifyRateLimit(config.rateLimit);
      returnConfig.rateLimit = newRateLimit;
    }
    return returnConfig;
  },
};

const Config = {
  mainTags: () => ConfigModel.getMainTags(),
  rateLimit: () => ConfigModel.getRateLimit(),
};

export default {
  Query,
  Mutation,
  Config,
};
