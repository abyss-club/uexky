import { ApolloError } from 'apollo-server-koa';

class AuthError extends ApolloError {
  constructor(message, ...params) {
    super(message, 'UNAUTHENTICATED', ...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AuthError);
    }
    this.name = 'AuthError';
    this.authError = true;
  }
}

class ParamsError extends ApolloError {
  constructor(message, ...params) {
    super(message, 'PARAMS_ERROR', ...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ParamsError);
    }
    this.name = 'ParamsError';
    this.paramsError = true;
  }
}

class InternalError extends ApolloError {
  constructor(message, ...params) {
    super(message, 'INTERNAL_ERROR', ...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, InternalError);
    }
    this.name = 'InternalError';
    this.internalError = true;
  }
}

class PermissionError extends ApolloError {
  constructor(message, ...params) {
    super(message, 'FORBIDDEN', ...params);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, PermissionError);
    }
    this.name = 'PermissionError';
    this.internalError = true;
  }
}

class NotFoundError extends ApolloError {
  constructor(message, ...params) {
    super(message, 'NOT_FOUND', ...params);
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
