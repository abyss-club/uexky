module.exports = {
  "extends": [
    "eslint:recommended",
    "plugin:node/recommended",
    "plugin:jest/recommended",
    "airbnb-base",
    "plugin:eslint-comments/recommended",
    "plugin:promise/recommended",
    "plugin:unicorn/recommended",
  ],
  "rules": {
    "no-use-before-define": [0, {}],
    "no-underscore-dangle": ["error", {"allow": ["__MONGO__", "__MONGO_URI__", "__MONGO_DB_NAME__", "_id"]}],
    "node/no-unsupported-features/es-syntax": 0,
    "unicorn/prevent-abbreviations": "off",
    "unicorn/filename-case": "off",
    "unicorn/catch-error-name": [
      "error",
      {
        "caughtErrorsIgnorePattern": "^.*Err$"
      }
    ],
  },
  "env": {
    "es6": true,
  },
  "globals": {
    "Atomics": "readonly",
    "SharedArrayBuffer": "readonly",
  },
  "parser": "babel-eslint",
  "parserOptions": {
    "sourceType": "module",
    "ecmaVersion": 2019,
  },
  "settings": {
    "import/resolver": {
      "alias": {
        "map": [
          ["~", "./src"]
        ],
        "extensions": [".js", ".jsx", ".json"],
      },
    },
  },
};
