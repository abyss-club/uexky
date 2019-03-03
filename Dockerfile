FROM node:11.3-stretch

ADD package.json /usr/app/
ADD yarn.lock /usr/app/

WORKDIR /usr/app
RUN yarn

ADD . /usr/app
RUN yarn build

ENV MONGODB_URI=
ENV API_DOMAIN=
ENV MAIN_DOMAIN=
ENV REDIS_URI=
ENV PROTO=
ENV PORT=
ENV MAILGUN_PRIVATE_KEY=
ENV MAILGUN_PUBLIC_KEY=
ENV MAILGUN_DOMAIN=
ENV MAILGUN_SENDER=

CMD [ "yarn", "start" ]
