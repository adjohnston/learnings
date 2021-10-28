import * as TE from 'fp-ts/TaskEither'
import * as RTE from 'fp-ts/ReaderTaskEither'
import * as R from 'fp-ts/Reader'
import * as RE from 'fp-ts/ReaderEither'
import * as E from 'fp-ts/Either'
import { pipe } from 'fp-ts/function'

interface User {
  id: number
  name: string
  email: string
}

interface Item {
  id: number
  name: string
  qty: number
  unitPrice: number
}

interface Basket {
  id: number
  userId: number
  items: Item[]
}

interface Config {
  code: 'gb' | 'no'
  name: 'United Kingdom' | 'Norway'
  vatRate: number
}

const gbConfig: Config = {
  code: 'gb',
  name: 'United Kingdom',
  vatRate: 0.2,
}

const noConfig: Config = {
  code: 'no',
  name: 'Norway',
  vatRate: 0.25,
}

// mock of the user DB.
const userDB = [
  { id: 1, name: 'Adam', email: 'adam@adam.co.uk' },
  { id: 2, name: 'Dave', email: 'dave@allthingsdave.com' },
  { id: 3, name: 'Laura', email: 'laura@python-emporium.biz' },
]

// mock the basket DB.
const basketDB = [
  {
    id: 1,
    userId: 1,
    items: [
      { id: 1, name: 'Puggo', qty: 1, unitPrice: 10000 },
      { id: 5, name: 'Dog food', qty: 3, unitPrice: 250 },
    ],
  },
  {
    id: 2,
    userId: 3,
    items: [
      { id: 2, name: 'Python', qty: 4, unitPrice: 3500 },
      { id: 3, name: 'Reptile house', qty: 2, unitPrice: 1000 },
      { id: 8, name: 'Mice', qty: 20, unitPrice: 500 },
    ],
  },
]

const getUser = (email: string): Promise<User> => {
  try {
    const user = userDB.find(user => user.email === email)
    if (!user) {
      throw Error('No user with given email')
    }
    return Promise.resolve(user)
  } catch (error: unknown) {
    const errorMessage =
      error instanceof Error ? error.message : 'Could not fetch user'
    return Promise.reject(errorMessage)
  }
}

const getBasket = (userId: number): Promise<Basket> => {
  try {
    const basket = basketDB.find(basket => basket.userId === userId)
    if (!basket) {
      throw Error('No basket found with given user ID')
    }
    return Promise.resolve(basket)
  } catch (error: unknown) {
    const errorMessage =
      error instanceof Error ? error.message : 'Could not fetch basket'
    return Promise.reject(errorMessage)
  }
}

const basketContainsPug = (items: Basket['items']) => (config: Config) =>
  config.code === 'no' && items.find(item => item.name === 'Puggo')

// version of validateBasket without lifting the basketContainsPug
// into a context of ReaderEither.
// const validateBasket = (
//   basket: Basket,
// ): RE.ReaderEither<Config, string, Basket> =>
//   pipe(
//     RE.ask<Config>(),
//     RE.chain<Config, string, Config, Basket>(config =>
//       basketContainsPug(basket.items)(config)
//         ? RE.left('We cannot sell Puggos in your country')
//         : RE.right(basket),
//     ),
//   )

const validateBasket = (
  basket: Basket,
): RE.ReaderEither<Config, string, Basket> =>
  pipe(
    RE.of(basketContainsPug),
    RE.ap(RE.of(basket.items)),
    RE.ap(RE.ask()),
    RE.chain(isValid =>
      isValid
        ? RE.left('We cannot sell Puggos in your country')
        : RE.right(basket),
    ),
  )

const calculateTotal = (basket: Basket): R.Reader<Config, number> =>
  pipe(
    R.ask<Config>(),
    R.map<Config, number>(config =>
      pipe(
        basket,
        basket => basket.items,
        items =>
          items.map(
            item => item.unitPrice + item.unitPrice * config.vatRate * item.qty,
          ),
        totals => totals.reduce((x, y) => x + y, 0),
      ),
    ),
  )

const calculateBasketTotal = (email: string): any =>
  pipe(
    // fetch a user so we can get their basket.
    TE.tryCatch<string, User>(
      () => getUser(email),
      error => error,
    ),
    // fetch a user's basket
    TE.chain<string, User, Basket>(user =>
      TE.tryCatch<string, Basket>(
        () => getBasket(user.id),
        error => error,
      ),
    ),
    // switch to from task either to reader task either so that we can
    // "inject" our config.
    RTE.fromTaskEither,
    // validateBasket is not a reader task either it's a reader either
    // so we use a Kleisli chain to.
    RTE.chainReaderEitherK(validateBasket),
    RTE.chainReaderK(calculateTotal),
  )

calculateBasketTotal('adam@adam.co.uk')(gbConfig)().then(
  // calculateBasketTotal('dave@allthingsdave.com')(noConfig)().then(
  // calculateBasketTotal('laura@python-emporium.biz')(noConfig)().then(
  E.fold(
    e => console.log('Got an error:', e),
    r => console.log('Total is:', r),
  ),
)
