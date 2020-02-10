const fastify = require('fastify')({ logger: true })
const routes = require('./routes')
const db = require('./db')

fastify.register(db, { url: 'mongodb://localhost:27017/' })
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
