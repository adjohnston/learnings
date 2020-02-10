const fastify = require('fastify')({ logger: true })
const routes = require('./routes')

fastify.register(routes)

const start = async () => {
  try {
    await fastify.listen(3000)
  } catch (error) {
    fastify.log.error(error)
    process.exit(1)
  }
}

start()
