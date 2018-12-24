import UserModel from '../models/user';

// const User = (user) => {
//   profile: (obj, args, ctx, info) => {
//     console.log('resolver ctx', ctx);
//     return UserModel.profile(ctx);
//   }
// };

const profile = (ctx) => {
  console.log(ctx);
  if (!ctx.user) return null;
  return ctx.user;
};

// export default User;
export { profile };
