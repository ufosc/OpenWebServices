const { version } = require('./package.json');

export const API_ENDPOINT = process.env.NEXT_PUBLIC_API_ENDPOINT || "http://localhost:8080"
export const VERSION = 'v'+version
