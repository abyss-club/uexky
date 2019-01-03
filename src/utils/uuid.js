import mongoose from 'mongoose';

const { ObjectId } = mongoose.Types;
const code = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_';

// Convert 3-hex number to 2-codes base64 string
const hex3ToBase64 = (hex) => {
  if (hex.length !== 3) {
    throw new Error(`invalid params: ${hex}`);
  }
  const padTo = 2;
  let remaining = parseInt(hex, 16);
  const codes = [];
  while (remaining > 0) {
    codes.push(code[remaining % 64]);
    remaining = Math.floor(remaining / 64);
  }
  while (codes.length < padTo) {
    codes.push('A');
  }
  return codes.reverse().join('');
};

// Convert 2-codes base64 string to 3-hex number
const base64ToHex3 = (base64) => {
  if (base64.length !== 2) {
    throw new Error(`invalid params: ${base64}`);
  }
  let num = 0;
  for (let i = 0; i < 2; i += 1) {
    const index = code.indexOf(base64[i]);
    if (index < 0) {
      throw new Error('Invalid uuid');
    }
    num = num * 64 + index;
  }
  return num.toString(16).padStart(3, '0');
};

const encode = (objectId) => {
  const idStr = objectId.valueOf();
  if (idStr.length !== 24) {
    throw new Error('Invalid objectId');
  }
  const hex = idStr.substring(9, 24) + idStr.substring(0, 9);
  const codes = [];
  for (let i = 0; i < hex.length; i += 3) {
    codes.push(hex3ToBase64(hex.substring(i, i + 3)));
  }
  return codes.join('');
};

const decode = (uuid) => {
  if (uuid.length !== 8) {
    throw new Error('Invalid uuid');
  }
  const ids = [];
  const base64 = uuid.substring(3, 8) + uuid.substring(0, 3);
  for (let i = 0; i < uuid.length; i += 2) {
    ids.push(base64ToHex3(base64.substring(i, i + 2)));
  }
  const idStr = ids.join('');
  return ObjectId(idStr);
};

const genRandomStr = (length) => {
  if (Number.isNaN(length) || length < 1) throw new Error('Invalid length');
  const str = [];
  const getRandomInt = () => Math.floor(Math.random() * Math.floor(code.length));
  for (let i = 0; i < length; i += 1) {
    str.push(code[getRandomInt()]);
  }
  return str.join('');
};

export { encode, decode, genRandomStr };
