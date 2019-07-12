import { ParamsError } from '~/utils/error';
import { query } from '~/utils/pg';

// type SliceQuery
// {
//   before: String
//   after: String
//   limit: Int
// }
// type option {
//   select: <select clause>,
//   where: <where clause>,
//   before: <function, return where clause deal with before>,
//   after: <function, return where clause deal with after>,
//   order: <order clause>,
//   desc: <bool, is need desc when after enable>,
//   params: <pg query params>,
//   name: <slice data field name>,
//   make: <function, wrap pg row>,
//   toCursor: <function, calculate cursor>,
// };
// return type sliceInfo {
//   [sliceName]: [result],
//   sliceInfo: {
//     firstCursor: '',
//     lastCursor: '',
//     hasNext: true,
//   }
// }

const enable = op => typeof op === 'string';


async function querySlice(sq, {
  select, where, before, after, order, desc, params, name, make, toCursor,
}) {
  if ((enable(sq.before) && enable(sq.after)) || (!enable(sq.before) && !enable(sq.after))) {
    throw new ParamsError('you must specific either before or after');
  }

  const needDesc = ((enable(sq.after) && desc) || (enable(sq.before) && !desc));
  // query
  const sql = [select, where];
  if (sq.before) {
    sql.push(where && 'AND');
    sql.push(before(sq.before));
  }
  if (sq.after) {
    sql.push(where && 'AND');
    sql.push(after(sq.after));
  }
  sql.push(order);
  if (needDesc) {
    sql.push('DESC');
  }
  sql.push(`LIMIT ${sq.limit + 1}`);
  const { rows } = await query(sql.join(' '), params);
  // parse result
  const len = (rows || []).length;
  const slice = {
    [name]: [],
    sliceInfo: { firstCursor: '', lastCursor: '', hasNext: false },
  };
  if (len === 0) {
    return slice;
  }
  if (desc !== needDesc) {
    rows.reverse();
  }
  if (len <= sq.limit) {
    slice[name] = rows.map(row => make(row));
  } else { // len === sq.limit + 1
    slice.sliceInfo.hasNext = true;
    if (enable(sq.before)) {
      slice[name] = rows.slice(1, len).map(row => make(row));
    } else {
      slice[name] = rows.slice(0, len - 1).map(row => make(row));
    }
  }
  slice.sliceInfo.firstCursor = toCursor(slice[name][0]);
  slice.sliceInfo.lastCursor = toCursor(slice[name][slice[name].length - 1]);
  return slice;
}

export default querySlice;
