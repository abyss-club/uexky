import mongoose from 'mongoose';
import { InternalError, ParamsError } from '~/error';

const ConfigSchema = new mongoose.Schema({
  optionName: { type: String, required: true, unique: true },
  optionValue: { type: String, required: true },
  autoload: { type: Boolean },
});

ConfigSchema.statics.getMainTags = async function getMainTags() {
  const result = await ConfigModel.findOne({ optionName: 'mainTags' }, 'optionValue').orFail(() => new InternalError('mainTags not yet set.'));
  return JSON.parse(result.optionValue);
};

ConfigSchema.statics.modifyMainTags = async function modifyMainTags(tags) {
  if (!Array.isArray(tags) || !tags.length) {
    throw new ParamsError('Provided tags is not a non-empty array.');
  }
  const newMainTags = { optionName: 'mainTags', optionValue: JSON.stringify(tags) };
  await ConfigModel.updateOne({ optionName: 'mainTags' }, newMainTags, { upsert: true });
};

const ConfigModel = mongoose.model('Config', ConfigSchema);

export default ConfigModel;
