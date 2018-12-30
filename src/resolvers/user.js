import { AuthenticationError } from 'apollo-server-koa';
import UserModel from '../models/user';
import AuthModel from '../models/auth';

const resolvers = {
  User: ({ email, name, tags }) => ({ email, name, tags }),

  Query: {
    profile: (_, __, ctx) => {
      if (!ctx.user) throw new AuthenticationError('Authentication needed.');
      const { email, name, tags } = ctx.user;
      return { email, name, tags };
    },
  },

  Mutation: {
    auth: (obj, args, ctx, info) => {
      AuthModel.addToAuth(args.email);
      return true;
    },
  },
};

// export default UserTypes;
// export { profile };
export default resolvers;
