import gleam/list

fn recurse(number: Int, binaries: List(Int)) -> List(Int) {
  case number {
    0 -> binaries
    _ as n -> recurse(number / 2, list.append(binaries, [n % 2]))
  }
}

pub fn egg_count(number: Int) -> Int {
  recurse(number, [])
  |> list.filter(fn(binary: Int) -> Bool { binary == 1 })
  |> list.length()
}
