import mongoose from 'mongoose';
import app from './app';

const dbHost = process.env.DB_HOST;
const dbPort = process.env.DB_PORT;
const dbName = process.env.DB_NAME;
mongoose.connect(`mongodb://${dbHost}:${dbPort}/${dbName}`, { useNewUrlParser: true });

const port = process.env.PORT || 5000;
app.listen(port);
console.info(`Listening to http://localhost:${port} 🚀`);
