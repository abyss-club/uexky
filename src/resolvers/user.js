import UserModel from '../models/user';
import AuthModel from '../models/auth';

const resolvers = {
  User: profile => profile,

  Query: {
    profile: (_, __, ctx) => {
      if (!ctx.user) return null;
      const { email, name, tags } = ctx.user;
      return { email, name, tags };
    },
  },
};

// export default UserTypes;
// export { profile };
export default resolvers;
