import { ifElse, isNil, complement, identity, tap } from 'ramda'
import * as TE from 'fp-ts/TaskEither'
import * as RTE from 'fp-ts/ReaderTaskEither'
import * as R from 'fp-ts/Reader'
import * as RE from 'fp-ts/ReaderEither'
import * as E from 'fp-ts/Either'
import { pipe } from 'fp-ts/function'

/**
 * ============================================================================
 * TYPES
 * ============================================================================
 */

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
  vatRate: number
}

/**
 * ============================================================================
 * CONFIG DEPENDENCIES
 * ============================================================================
 */

const gbConfig: Config = {
  code: 'gb',
  vatRate: 0.2,
}

const noConfig: Config = {
  code: 'no',
  vatRate: 0.25,
}

/**
 * ============================================================================
 * DB MOCKS
 * ============================================================================
 *
 * mock of the user DB.
 */
const userDB = [
  { id: 1, name: 'Adam', email: 'adam@adam.co.uk' },
  { id: 2, name: 'Dave', email: 'dave@allthingsdave.com' },
  { id: 3, name: 'Laura', email: 'laura@python-emporium.biz' },
]

/**
 * mock of the basket DB.
 */
const basketDB = [
  {
    id: 1,
    userId: 1,
    items: [
      { id: 1, name: 'Zebraffidile', qty: 1, unitPrice: 10000 },
      { id: 5, name: 'Fresh leaves', qty: 3, unitPrice: 250 },
    ],
  },
  {
    id: 2,
    userId: 3,
    items: [
      { id: 2, name: 'Snek', qty: 4, unitPrice: 3500 },
      { id: 3, name: 'Reptile house', qty: 2, unitPrice: 1000 },
      { id: 8, name: 'Mice', qty: 20, unitPrice: 500 },
    ],
  },
]

/**
 * ============================================================================
 * HELPERS
 * ============================================================================
 */

const getUser = (email: string): Promise<User> => {
  try {
    const user = userDB.find(user => user.email === email)

    if (!user) {
      throw Error('No user with given email')
    }

    return Promise.resolve(user)
  } catch (error: unknown) {
    const response =
      error instanceof Error ? error : new Error('Could not fetch user')

    return Promise.reject(response)
  }
}

const getBasket = (userId: number): Promise<Basket> => {
  try {
    const basket = basketDB.find(
      (basket: Basket): boolean => basket.userId === userId,
    )

    if (!basket) {
      throw Error('No basket found with given user ID')
    }

    return Promise.resolve(basket)
  } catch (error: unknown) {
    const response =
      error instanceof Error ? error : new Error('Could not fetch basket')

    return Promise.reject(response)
  }
}

const basketContainsImaginaryAnimal =
  (items: Basket['items']) =>
  (config: Config): boolean =>
    config.code === 'no' &&
    Boolean(items.filter(item => item.name === 'Zebraffidile').length)

// version of validateBasket without lifting the basketContainsImaginaryAnimal
// into a context of ReaderEither.
// const validateBasket = (
//   basket: Basket,
// ): RE.ReaderEither<Config, string, Basket> =>
//   pipe(
//     RE.ask<Config>(),
//     RE.chain<Config, string, Config, Basket>(config =>
//       basketContainsImaginaryAnimal(basket.items)(config)
//         ? RE.left('We cannot sell Zebraffidiles in your country')
//         : RE.right(basket),
//     ),
//   )

const validateBasket = (
  basket: Basket,
): RE.ReaderEither<Config, Error, Basket> =>
  pipe(
    RE.of(basketContainsImaginaryAnimal(basket.items)),
    RE.ap(RE.ask()),
    RE.chain(
      (isValid: boolean): RE.ReaderEither<Config, Error, Basket> =>
        isValid
          ? RE.left(new Error('We cannot sell Zebraffidiles in your country'))
          : RE.right(basket),
    ),
  )

const calculateTotal = (basket: Basket): R.Reader<Config, number> =>
  pipe(
    R.ask<Config>(),
    R.map<Config, number>(({ vatRate }) =>
      pipe(
        basket,
        (basket: Basket): Basket['items'] => basket.items,
        (items: Basket['items']): number[] =>
          items.map(
            ({ unitPrice, qty }): number =>
              (unitPrice + unitPrice * vatRate) * qty,
          ),
        (totals: number[]): number =>
          totals.reduce((x: number, y: number): number => x + y, 0),
      ),
    ),
  )

const createUnknownError = (): Error => Error('Unknown error')

const calculateBasketTotal = (
  email: string,
): RTE.ReaderTaskEither<Config, Error, number> =>
  pipe(
    /**
     * fetch a user so we can get their basket.
     */
    TE.tryCatch<Error, User>(
      (): Promise<User> => getUser(email),
      ifElse(complement(isNil), identity, createUnknownError),
    ),
    /**
     * fetch a user's basket.
     */
    TE.chain<Error, User, Basket>(
      (user: User): TE.TaskEither<Error, Basket> =>
        TE.tryCatch<Error, Basket>(
          (): Promise<Basket> => getBasket(user.id),
          ifElse(complement(isNil), identity, createUnknownError),
        ),
    ),
    /**
     * switch to from task either to reader task either so that we can "inject"
     * our config.
     */
    RTE.fromTaskEither,
    /**
     * validateBasket is not a reader task either it's a reader either so we
     * use a Kleisli chain to call it.
     */
    RTE.chainReaderEitherK(validateBasket),
    RTE.chainReaderK(calculateTotal),
  )

/**
 * ============================================================================
 * MAIN
 * ============================================================================
 */

const main = calculateBasketTotal('adam@adam.co.uk')(gbConfig)

main().then(
  // calculateBasketTotal('dave@allthingsdave.com')(gbConfig)().then(
  // calculateBasketTotal('laura@python-emporium.biz')(gbConfig)().then(
  E.fold(
    (e: Error): void => console.log(e.message),
    (r: number): void => console.log('Total is:', r),
  ),
)
