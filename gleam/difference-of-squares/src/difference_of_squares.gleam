import gleam/list
import gleam/result

pub fn square_of_sum(n: Int) -> Int {
  let sum =
    list.range(1, n)
    |> list.reduce(fn(sum, x) { sum + x })
    |> result.unwrap(0)

  sum * sum
}

pub fn sum_of_squares(n: Int) -> Int {
  list.range(1, n)
  |> list.map(fn(x) { x * x })
  |> list.reduce(fn(sum, x) { sum + x })
  |> result.unwrap(0)
}

pub fn difference(n: Int) -> Int {
  square_of_sum(n) - sum_of_squares(n)
}
