class AuthError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AuthError);
    }
    this.authError = true;
  }
}

class ParamsError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ParamsError);
    }
    this.paramsError = true;
  }
}

class InternalError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, InternalError);
    }
    this.internalError = true;
  }
}

class PermissionError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, PermissionError);
    }
    this.internalError = true;
  }
}

export {
  AuthError, ParamsError, InternalError, PermissionError,
};
