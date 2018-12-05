import { gql } from 'apollo-server-koa';

import { base } from './base';
import { notificationQueries, notificationTypes } from './notification';
import { postMutations, postQueries, postTypes } from './post';
import { tagQueries, tagTypes } from './tag';
import { threadQueries } from './thread';
import { userMutations, userTypes } from './user';

export default gql`${base}
${notificationTypes}
${postTypes}
${tagTypes}
${userTypes}
  type Query {
    ${notificationQueries}
    ${postQueries}
    ${tagQueries}
    ${threadQueries}
  }
  type Mutation {
    ${postMutations}
    ${userMutations}
  }
`;
