const { GraphQLScalarType, Kind } = require('graphql');

const Time = new GraphQLScalarType({
  name: 'Time',
  description: 'Time scalar type, data ISO string',
  parseValue(value) {
    return new Date(value);
  },
  serialize(value) {
    return value.toISOString();
  },
  parseLiteral(ast) {
    if (ast.kind === Kind.String) {
      return new Date(ast.value);
    }
    return null;
  },
});

// Default Types Resolvers:
//   SliceInfo:
//     firstCursor, lastCursor, hasNext


export default {
  Time,
};
