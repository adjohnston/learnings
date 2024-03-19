import gleam/list

pub fn square_of_sum(n: Int) -> Int {
  list.range(1, n)
  |> list.fold(0, fn(sum, x) { sum + x })
  |> fn(x) { x * x }
}

pub fn sum_of_squares(n: Int) -> Int {
  list.range(1, n)
  |> list.fold(0, fn(sum, x) { sum + x * x })
}

pub fn difference(n: Int) -> Int {
  square_of_sum(n) - sum_of_squares(n)
}
