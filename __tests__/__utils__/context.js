import UserModel from '~/models/user';
import { query } from '~/utils/pg';

async function mockContext({ email, name, role }) {
  let auth = await UserModel.authContext({ email });
  if ((name || '') === '' && (role || '') === '') {
    return { auth };
  }

  const user = auth.signedInUser();
  if (name) {
    await query('UPDATE public.user SET name=$1 WHERE id=$2', [name, user.id]);
  }
  if (role) {
    await query('UPDATE public.user SET role=$1 WHERE id=$2', [role, user.id]);
  }
  auth = await UserModel.authContext({ email });
  return { auth };
}

export default mockContext;
