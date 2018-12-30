import { merge } from 'lodash';
import UserResolver from './user';
import TagResolver from './tag';

// const resolvers = Object.assign({}, UserResolver, TagResolver);
const resolvers = merge(UserResolver, TagResolver);
export default resolvers;
