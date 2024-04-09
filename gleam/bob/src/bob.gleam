import gleam/string

pub fn hey(remark: String) -> String {
  let trimmed_remark = string.trim(remark)
  let is_question = string.ends_with(trimmed_remark, "?")
  let is_empty = string.is_empty(trimmed_remark)

  let is_shout =
    string.uppercase(trimmed_remark) == trimmed_remark
    && string.lowercase(trimmed_remark) != trimmed_remark

  case trimmed_remark {
    _ if is_empty -> "Fine. Be that way!"
    _ if is_question && is_shout -> "Calm down, I know what I'm doing!"
    _ if is_question -> "Sure."
    _ if is_shout -> "Whoa, chill out!"
    _ -> "Whatever."
  }
}
