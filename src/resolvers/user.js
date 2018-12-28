import UserModel from '../models/user';
import AuthModel from '../models/auth';

const resolvers = {
  User: ({ email, name, tags }) => ({ email, name, tags }),

  Query: {
    profile: (_, __, ctx) => {
      if (!ctx.user) return {};
      const { email, name, tags } = ctx.user;
      return { email, name, tags };
    },
  },
};

// export default UserTypes;
// export { profile };
export default resolvers;
