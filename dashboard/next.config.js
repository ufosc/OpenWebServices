const { version } = require('./package.json')

/** @type {import('next').NextConfig} */
module.exports = {
  "output": 'standalone',
  env: { version },
}
