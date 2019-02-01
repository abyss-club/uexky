import mongoose from 'mongoose';

const RateLimiterSchema = new mongoose.Schema({
  ip: String,
  email: String,
  mutable: Boolean,
  updatedAt: Date,
  remaining: Number,
});
const RateLimiterModel = mongoose.model('rate_limit', RateLimiterSchema);

// TODO: add ttl for limiter
// TODO: cut down limiters' remaining
// TODO: read settings from config

export default RateLimiterModel;
