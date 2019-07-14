import { connectDb, query } from '~/utils/pg';
import UserModel from '~/models/user';
import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';
import env from '~/utils/env';

const startPg = async () => {
  const pgPool = await connectDb(env.PG_URI);
  return pgPool;
};

function randInt(min, max) {
  const a = Math.ceil(min);
  const b = Math.floor(max);
  return Math.floor(Math.random() * (b - a)) + a;
}

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

const makeDemo = async () => {
  await startPg();

  // insert user
  const ctx0 = await mockContext({ email: 't1@abyss.club', name: 'L0' });
  const ctx1 = await mockContext({ email: 't2@abyss.club', name: 'L1' });
  const ctx2 = await mockContext({ email: 't3@abyss.club', name: 'L2' });
  const ctx3 = await mockContext({ email: 't4@abyss.club', name: 'L3', role: 'mod' });
  const ctx4 = await mockContext({ email: 't5@abyss.club', name: 'L4', role: 'admin' });
  const ctxList = [ctx0, ctx1, ctx2, ctx3, ctx4];
  console.log(ctxList);

  // mock tags
  const mainTags = ['MainA', 'MainB', 'MainC'];
  await Promise.all(mainTags.map(mt => query(
    'INSERT INTO tag (name, is_main) VALUES ($1, $2)',
    [mt, true],
  )));

  // pub thread
  /* eslint-disable no-await-in-loop */
  const threadIds = [];
  for (let i = 0; i < 20; i += 1) {
    const subTagsCount = randInt(0, 5);
    const subTags = [];
    for (let j = 0; j < subTagsCount; j += 1) {
      subTags.push(`Sub${randInt(0, 100)}`);
    }
    const ctx = ctxList[randInt(0, 5)];
    const anonymous = randInt(0, 2) === 1;
    const params = {
      ctx,
      thread: {
        anonymous,
        title: `thread${i}`,
        content: `content${i}`,
        mainTag: mainTags[randInt(0, 3)],
        subTags: [...new Set(subTags)],
      },
    };
    console.log('new thread', params);
    const thread = await ThreadModel.new(params);
    threadIds.push(thread.id);
  }

  const postIds = [];
  for (let i = 0; i < 1000; i += 1) {
    const ctx = ctxList[randInt(0, 5)];
    const threadId = threadIds[randInt(0, 10)];
    const anonymous = randInt(0, 2) === 1;
    const quoteIds = [];
    const quoteCount = randInt(0, 4);
    if (postIds.length > 2) {
      for (let j = 0; j < quoteCount; j += 1) {
        quoteIds.push(postIds[randInt(0, postIds.length)]);
      }
    }
    const params = {
      ctx,
      post: {
        threadId,
        anonymous,
        content: `post${i}`,
        quoteIds: [...new Set(quoteIds)],
      },
    };
    console.log(params);
    const post = await PostModel.new(params);
    postIds.push(post.id);
  }
  /* eslint-enable no-await-in-loop */
};

(async () => {
  try {
    await makeDemo();
  } catch (e) {
    console.log(e.stack);
  }
})();
