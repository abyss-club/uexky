import { ensureSignIn } from '../models/user';
import AuthModel from '../models/auth';

const Query = {
  profile: (_, __, ctx) => ensureSignIn(ctx),
};

const Mutation = {
  auth: (obj, { email }, ctx) => {
    if (ctx.user) throw new Error('already signed in!');
    AuthModel.addToAuth(email);
    return true;
  },
  setName: (obj, { name }, ctx) => {
    const user = ensureSignIn(ctx);
    user.setName(name);
  },
  syncTags: (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    user.syncTags(tags);
  },
  addSubbedTags: (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    user.addSubbedTags(tags);
  },
  delSubbedTags: (obj, { tags }, ctx) => {
    const user = ensureSignIn(ctx);
    user.delSubbedTags(tags);
  },
};

// Default Types Resolver:
//   User:
//     email, name, tags

export default {
  Query,
  Mutation,
};
