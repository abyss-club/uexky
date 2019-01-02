class AuthFail extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AuthFail);
    }
    this.authFail = true;
  }
}

export default AuthFail;
