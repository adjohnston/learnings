const fastifyPlugin = require('fastify-plugin')
const MongoClient = require('mongodb').MongoClient

const db = async (fastify, { url, ...options }) => {
  const db = await MongoClient.connect(url, options)

  fastify.decorate('mongo', db)
}

module.exports = fastifyPlugin(db)
