import * as mongoose from 'mongoose';
import app from './app';

const { DB_PORT, DB_HOST, DB_NAME } = process.env;
mongoose.connect(`mongodb://${DB_HOST}:${DB_PORT}/${DB_NAME}`, { useNewUrlParser: true });

const port = process.env.PORT || 5000;
const server = app.listen(port);
console.info(`Listening to http://localhost:${port} 🚀`);
