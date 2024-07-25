import gleam/option.{type Option}

pub type Player {
  Player(name: Option(String), level: Int, health: Int, mana: Option(Int))
}

pub fn introduce(player: Player) -> String {
  case player.name {
    option.Some(name) -> name
    option.None -> "Mighty Magician"
  }
}

pub fn revive(player: Player) -> Option(Player) {
  case player.health {
    0 -> {
      case player.level {
        lvl if lvl >= 10 ->
          option.Some(Player(..player, health: 100, mana: option.Some(100)))
        _ -> option.Some(Player(..player, health: 100))
      }
    }
    _ -> option.None
  }
}

pub fn cast_spell(player: Player, cost: Int) -> #(Player, Int) {
  case player.mana {
    option.Some(m) if m >= cost -> #(
      Player(
        ..player,
        mana: option.then(player.mana, fn(m) { option.Some(m - cost) }),
      ),
      cost * 2,
    )
    option.Some(m) if m < cost -> #(player, 0)
    _ -> {
      let remaining_health = player.health - cost

      case remaining_health {
        x if x <= 0 -> #(Player(..player, health: 0), 0)
        _ -> #(Player(..player, health: remaining_health), 0)
      }
    }
  }
}
