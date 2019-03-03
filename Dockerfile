FROM node:11.3-stretch

ENV PORT 3000

ADD package.json /usr/app/
ADD yarn.lock /usr/app/

WORKDIR /usr/app
RUN yarn

ADD . /usr/app
RUN yarn build

CMD [ "yarn", "start" ]
