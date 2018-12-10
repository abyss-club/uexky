module.exports = {
  "extends": ["eslint:recommended", "plugin:node/recommended", "airbnb-base"],
  "rules": {
    "no-use-before-define": [0, {}],
    "no-underscore-dangle": ["error", {"allow": ["__MONGO__"]}],
    "strict": 0,
  },
  "settings": {
  },
  "env": {
    "es6": true,
  },
  "parser": "babel-eslint",
  "parserOptions": {
    "sourceType": "module",
    "ecmaVersion": 2018,
  }
};
