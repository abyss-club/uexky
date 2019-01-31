import Generator from './generator';

const Base64 = {
  code: 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_',
  // Convert 2-codes base64 string to 3-hex number
  convertTohex3(b64Str) {
    if (b64Str.length !== 2) {
      throw new Error(`invalid params: ${b64Str}`);
    }
    let num = 0;
    for (let i = 0; i < 2; i += 1) {
      const index = this.code.indexOf(b64Str[i]);
      if (index < 0) {
        throw new Error('Invalid uuid');
      }
      num = num * 64 + index;
    }
    return num.toString(16).padStart(3, '0');
  },
  // Convert 3-hex number to 2-codes base64 string
  parseFromHex3(hexStr) {
    if (hexStr.length !== 3) {
      throw new Error(`invalid params: ${hexStr}`);
    }
    const padTo = 2;
    let remaining = parseInt(hexStr, 16);
    const codes = [];
    while (remaining > 0) {
      codes.push(this.code[remaining % 64]);
      remaining = Math.floor(remaining / 64);
    }
    while (codes.length < padTo) {
      codes.push('A');
    }
    return codes.reverse().join('');
  },
  randomString(len) {
    if (Number.isNaN(len) || len < 1) throw new Error('Invalid length');
    const str = [];
    const getRandomInt = () => Math.floor(Math.random() * this.code.length);
    for (let i = 0; i < len; i += 1) {
      str.push(this.code[getRandomInt()]);
    }
    return str.join('');
  },
};

const hexReverse = {
  0: '0', // 0000 - 0000
  1: '8', // 0001 - 1000
  2: '4', // 0010 - 0100
  3: 'c', // 0011 - 1100
  4: '2', // 0100 - 0010
  5: 'a', // 0101 - 1010
  6: '6', // 0110 - 0110
  7: 'e', // 0111 - 1110
  8: '1', // 1000 - 0001
  9: '9', // 1001 - 1001
  a: '5', // 1010 - 0101
  b: 'd', // 1011 - 1101
  c: '3', // 1100 - 0011
  d: 'b', // 1101 - 1011
  e: '7', // 1110 - 0111
  f: 'f', // 1111 - 1111
};

const reverseTimestamp = (hex) => {
  const chars = [];
  for (let i = 0; i < hex.length; i += 1) {
    chars.push(hexReverse[hex[i]]);
  }
  return chars.reverse().join();
};

const UID = {
  async New() {
    const newId = await Generator.New();
    return newId;
  },
  encode(displayId) {
    if (displayId.length !== 10) {
      throw new Error(`Invalid uid display: ${displayId}`);
    }
    const hexes = [];
    for (let i = 0; i < displayId.length; i += 2) {
      hexes.push(Base64.convertTohex3(displayId.substring(i, i + 2)));
    }
    const hex = hexes.join('');
    return reverseTimestamp(hex.substring(0, 8)) + hex.substring(8, 15);
  },
  decode(storageId) {
    if (storageId.length !== 15) {
      throw new Error(`Invalid storaged uid: ${storageId}`);
    }
    // reverse timestamp
    const hex = reverseTimestamp(storageId.substring(0, 8)) + storageId.substring(8, 15);
    const codes = [];
    for (let i = 0; i < hex.length; i += 3) {
      codes.push(Base64.parseFromHex3(hex.substring(i, i + 3)));
    }
    return codes.join('');
  },
};

export default UID;
export { Base64 };
