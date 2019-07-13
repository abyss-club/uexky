const { GraphQLScalarType, Kind } = require('graphql');

const Time = new GraphQLScalarType({
  name: 'Time',
  description: 'Date custom scalar type',
  parseValue(value) {
    return new Date(value); // value from the client
  },
  serialize(value) {
    return value.getTime(); // value sent to the client
  },
  parseLiteral(ast) {
    if (ast.kind === Kind.INT) {
      return new Date(ast.value); // ast value is always in string format
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
