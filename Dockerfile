FROM node:latest

WORKDIR /app
ADD package.json ./
ADD yarn.lock ./

RUN yarn

ADD . ./
RUN yarn build

ENV PORT=
ENV MONGODB_URI=
ENV API_DOMAIN=
ENV MAIN_DOMAIN=
ENV REDIS_URI=
ENV PROTO=
ENV MAILGUN_PRIVATE_KEY=
ENV MAILGUN_PUBLIC_KEY=
ENV MAILGUN_DOMAIN=
ENV MAILGUN_SENDER=

EXPOSE 5000

CMD [ "yarn", "start" ]
