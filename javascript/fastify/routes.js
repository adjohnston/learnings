const routes = async (fastify, options) => {
  const db = fastify.mongo.db('db')
  const collection = db.collection('test')

  fastify.get('/', async (request, response) => {
    return { hello: 'world' }
  })

  fastify.get('/search/:id', async (request, response) => {
    const result = await collection.findOne({ id: request.params.id })

    if (result.value === null) {
      throw new Error('No value found')
    }

    return result.value
  })
}

module.exports = routes
