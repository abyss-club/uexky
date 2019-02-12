function unicodeLength(str) {
  let len = 0;
  for (let i = 0; i < str.length; i += 1) {
    const w = (str.codePointAt(i) > 255) ? 2 : 1;
    len += w;
  }
  return len;
}

export default {
  unicodeLength,
};
