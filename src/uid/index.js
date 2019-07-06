import { ParamsError } from '~/utils/error';
import newSuid from './generator';

const Base64 = {
  code: 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_',
  // Convert 2-codes base64 string to 3-hex number
  convertToBigInt(b64Str) {
    let num = BigInt(0);
    for (let i = 0; i < b64Str.length; i += 1) {
      const n = this.code.indexOf(b64Str[i]);
      num = num * BigInt(64) + BigInt(n);
    }
    return num;
  },
  parseFromBigInt(bint, padTo) {
    const codes = [];
    const base = BigInt(64);
    let num = BigInt(bint);
    while (num > 0) {
      const r = num % base;
      codes.push(this.code[r]);
      num /= base;
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

const display2storage = (duid) => {
  if (typeof duid !== 'string' || duid.length < 6 || duid.length > 10) {
    throw new Error(`Invalid uid display: ${duid}`);
  }
  const len = duid.length;
  const raw = [duid.substring(1, len), duid[0]];
  return Base64.convertToBigInt(raw.join(''));
};

const storage2display = (suid) => {
  if (typeof suid !== typeof BigInt(0) || suid < BigInt(2 ** 28) || suid > BigInt(2 ** 60)) {
    throw new Error(`Invalid storaged uid: ${suid}`);
  }
  const raw = Base64.parseFromBigInt(suid);
  const len = raw.length;
  return [raw[len - 1], raw.substring(0, len - 1)].join('');
};

const UID = {
  parse(input) {
    if (typeof input === 'string') {
      return {
        duid: input,
        suid: display2storage(input),
        type: 'uid',
      };
    } if (typeof input === typeof BigInt(0)) {
      return {
        suid: input,
        duid: storage2display(input),
        type: 'uid',
      };
    } if (typeof input === 'object' && input.type === 'uid') {
      return input;
    }
    throw ParamsError(`unknown value: ${input}`);
  },
  async new() {
    const suid = await newSuid();
    return this.parse(suid);
  },
};

export default UID;
export { Base64 };
