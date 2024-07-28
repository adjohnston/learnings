import gleam/option
import gleam/regex

pub fn is_valid_line(line: String) -> Bool {
  case line {
    "[DEBUG]" <> _ | "[INFO]" <> _ | "[WARNING]" <> _ | "[ERROR]" <> _ -> True
    _ -> False
  }
}

pub fn split_line(line: String) -> List(String) {
  let assert Ok(regex) = regex.from_string("<[*-~=]*>")

  regex.split(regex, line)
}

pub fn tag_with_user_name(line: String) -> String {
  let assert Ok(regex) = regex.from_string("User\\s+(\\S+)")

  case regex.scan(regex, line) {
    [regex.Match(submatches: [option.Some(name)], ..)] ->
      "[USER] " <> name <> " " <> line
    _ -> line
  }
}
