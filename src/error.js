class AuthError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AuthError);
    }
    this.authError = true;
  }
}

class DataError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, DataError);
    }
    this.dataError = true;
  }
}

export { AuthError, DataError };
