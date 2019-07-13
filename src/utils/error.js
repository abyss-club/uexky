class AuthError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AuthError);
    }
    this.name = 'AuthError';
    this.authError = true;
  }
}

class ParamsError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ParamsError);
    }
    this.name = 'ParamsError';
    this.paramsError = true;
  }
}

class InternalError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, InternalError);
    }
    this.name = 'InternalError';
    this.internalError = true;
  }
}

class PermissionError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, PermissionError);
    }
    this.name = 'PermissionError';
    this.internalError = true;
  }
}

class NotFoundError extends Error {
  constructor(...params) {
    super(...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, NotFoundError);
    }
    this.name = 'NotFoundError';
    this.notFoundError = true;
  }
}

export {
  AuthError, ParamsError, InternalError, PermissionError, NotFoundError,
};
