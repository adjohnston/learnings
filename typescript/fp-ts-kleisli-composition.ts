import * as O from "fp-ts/Option";
import * as E from "fp-ts/Either";
import { pipe } from "fp-ts/function";

const add =
  (x: number) =>
  (y: number): number =>
    x + y;
const mul =
  (x: number) =>
  (y: number): number =>
    x * y;

pipe(5, add(5), mul(2));

const addO =
  (x: number) =>
  (y: number): O.Option<number> => {
    const result = x + y;
    return result ? O.some(result) : O.none;
  };

const mulO =
  (x: number) =>
  (y: number): O.Option<number> => {
    const result = x * y;
    return result ? O.some(result) : O.none;
  };

pipe(5, addO(5), mulO(2));

pipe(5, addO(5), O.chain(mulO(2)));

const addE =
  (x: number) =>
  (y: number): E.Either<number, number> => {
    const result = x + y;
    return result >= 1 ? E.right(result) : E.left(0);
  };

pipe(1, addE(0), E.chainOptionK(() => "Cannot multiply by 0")(mulO(0)));
