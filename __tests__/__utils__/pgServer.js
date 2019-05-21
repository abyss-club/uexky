import { connectDb } from '~/utils/pg';

const startPg = async () => {
  const pgPool = await connectDb('postgresql://localhost/test_uexky');
  return pgPool;
};

export default startPg;
