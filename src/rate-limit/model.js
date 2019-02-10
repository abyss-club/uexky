import mongoose from 'mongoose';

const RateLimiterSchema = new mongoose.Schema({
  ip: String,
  email: String,
  mutable: Boolean,
  createdAt: Date,
  remaining: Number,
});
const RateLimiterModel = mongoose.model('rate_limit', RateLimiterSchema);

// TODO: add ttl for limiter
// TODO: cut down limiters' remaining
// TODO: read settings from config
//
RateLimiterSchema.statics.take = async function take(
  ctx, cost, mutable, ip, email = '',
) {
  const selector = { ip, mutable };
  if (email !== '') {
    selector.email = email;
  }
  const config = await ctx.config.getRateLimit();
  const now = new Date();
  const remainingInit = mutable ? config.MutLimit : config.QueryLimit;
  const rl = await RateLimiterModel.findOneAndUpdate(
    selector,
    {
      $setOnInsert: { ...selector, createdAt: now, remaining: remainingInit },
      $dec: { remaining: cost },
    },
    { new: true, upsert: 1 },
  ).exec();
  if (rl.remaining < 0) {
    throw new Error('rate limit execceed!');
  }
};

export default RateLimiterModel;
