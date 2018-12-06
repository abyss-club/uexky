import * as request from 'supertest';
import app from '../src/app';

test('Hello world works', async () => {
    const response = await request(app.callback()).get('/');
    expect(response.status).toBe(200);
    expect(response.text).toBe('Hello World!');
});
