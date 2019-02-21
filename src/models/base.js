import { ParamsError } from '~/utils/error';

// type SliceQuery
// {
//   before: String
//   after: String
//   limit: Int
// }
// type option {
//   desc: Boolean[false],
//   field: String[_id],
//   sliceName: String[slice],
//   parse: Func[(value) => (value)],
//   toCursor: Func[(value) => (value || '')],
// }
// return type sliceInfo {
//   [sliceName]: [result],
//   sliceInfo: {
//     firstCursor: '',
//     lastCursor: '',
//     hasNext: true,
//   }
// }
async function findSlice(sliceQuery, model, option) {
  let { before, after } = sliceQuery;
  const { limit } = sliceQuery;
  if ((typeof before === 'undefined') && (typeof after === 'undefined')) {
    throw new ParamsError('Both before and after fields are empty');
  } else if (!Number.isInteger(limit) || limit < 1) {
    throw new ParamsError('Limit must be greater than 0.');
  }
  if (option.desc) {
    [before, after] = [after, before];
  }

  const field = option.field || '_id';
  const sliceName = option.sliceName || 'slice';
  const parse = option.parse || (value => value);
  const toCursor = (value) => {
    if (!value) {
      return '';
    }
    if (option.toCursor) {
      return option.toCursor(value);
    }
    return value;
  };
  const query = option.query || {};

  if (typeof before !== 'undefined') {
    if (before !== '') {
      query[field] = { $lt: parse(before) };
    }
  } else if (after !== '') {
    query[field] = { $gt: parse(after) };
  }

  const slice = await model.find(query, null, {
    limit: (limit + 1),
    sort: { [field]: option.desc ? -1 : 1 },
  });
  if (option.desc) {
    slice.reverse();
  }

  return {
    [sliceName]: slice,
    sliceInfo: {
      firstCursor: toCursor((slice[0] || {})[field]),
      lastCursor: toCursor((slice[slice.length - 1])[field]),
      hasNext: (slice.length) > limit,
    },
  };
}

export default findSlice;
