pub fn is_leap_year(year: Int) -> Bool {
  let is_4_year_leap = year % 4 == 0
  let is_100_year_leap = year % 100 == 0
  let is_400_year_leap = year % 400 == 0

  is_4_year_leap && { !is_100_year_leap || is_400_year_leap }
}
