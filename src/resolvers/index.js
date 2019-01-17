import { merge } from 'lodash';

import Base from './base';
import Config from './config';
import Notification from './notification';
import Post from './post';
import Tag from './tag';
import Thread from './thread';
import User from './user';

const resolvers = merge(Base, Config, Notification, Post, Tag, Thread, User);
export default resolvers;
