module.exports = {
  "extends": ["eslint:recommended", "plugin:node/recommended", "plugin:jest/recommended", "airbnb-base"],
  "rules": {
    "no-use-before-define": [0, {}],
    "no-underscore-dangle": ["error", {"allow": ["__MONGO__", "_id"]}],
    "strict": 0,
    "node/no-unsupported-features/es-syntax": 0,
  },
  "env": {
    "es6": true,
  },
  "parser": "babel-eslint",
  "parserOptions": {
    "sourceType": "module",
    "ecmaVersion": 2018,
  },
  settings: {
    'import/resolver': {
      alias: {
        map: [
          ['~', '.']
        ],
        extensions: ['.js', '.jsx', '.json'],
      },
    },
  },
};
