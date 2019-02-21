// import mongoose from 'mongoose';

// import app from '~/app';
import log from '~/utils/log';

// for babel-plugin-inline-dotenv
// const dbHost = process.env.DB_HOST;
// const dbPort = process.env.DB_PORT;
// const dbName = process.env.DB_NAME;

// mongoose.connect(`mongodb://${dbHost}:${dbPort}/${dbName}`, { useNewUrlParser: true });

const port = process.env.PORT || 5000;
// app.listen(port);

log.info(`Listening to http://localhost:${port} 🚀`);
