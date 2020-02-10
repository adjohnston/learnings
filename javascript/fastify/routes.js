const routes = async (fastify, options) => {
  fastify.get('/', async (request, response) => {
    return { hello: 'world' }
  })
}

module.exports = routes
