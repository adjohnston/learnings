module Main exposing (..)

import Browser
import Html exposing (Html, Attribute, div, input, text)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)

main : Program () Model Message
main =
    Browser.sandbox
        { init = init
        , view = view
        , update = update
        }

type alias Model = { content : String }

init : Model
init = { content = "" }

type Message = Change String

update : Message -> Model -> Model
update message model = 
  case message of
    Change nextMessage -> { model | content = nextMessage }

view : Model -> Html Message
view model = 
  div []
    [ input [ placeholder "Text to reverse", value model.content, onInput Change ] []
    , div [] [ text (String.reverse model.content) ] 
    ]