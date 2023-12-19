/** @type {import('next').NextConfig} */
module.exports = {
  async rewrites() {
      return [
      {
        source: '/auth/signup',
        destination: 'http://localhost:8080/auth/signup',
      },
      {
        source: '/auth/signin',
        destination: 'http://localhost:8080/auth/signin',
      },
      {
        source: '/auth/signout',
        destination: 'http://localhost:8080/auth/signout',
      },
      {
        source: '/auth/verify',
        destination: 'http://localhost:8080/auth/verify',
      },
      {
        source: '/auth/authorize',
        destination: 'http://localhost:8080/auth/authorize',
      },
      {
        source: '/auth/token',
        destination: 'http://localhost:8080/auth/token',
      },
      {
        source: '/client/:id',
        destination: 'http://localhost:8080/client/:id',
      },
      {
        source: '/client',
        destination: 'http://localhost:8080/client',
      },
      {
        source: '/user',
        destination: 'http://localhost:8080/user',
      },
    ]
  },
}
