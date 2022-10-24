---	CMPS4191 - Quiz #3
---	Rene Sanchez - 2018118383

module Todo exposing (Item)

import Html exposing (..)
import Html.Attributes exposing (..)

type alias Item =
    { name : String
    , details : String
    , priority : String
    , status: String
    }

initialModel : Item
initialModel =
    { name = ""
    , details = ""
    , priority = ""
    , status = ""
    }

view : Item -> Html msg
view item =
    div [ ]
        [ h1 [] [ text "ToDo" ]
        , Html.form []
            [ div []
                [ text "Name"
                , input [ id "name", type_ "text" ] []
                ]
            , div []
                [ text "Details"
                , input [ id "details", type_ "text" ] []
                ]
            , div []
                [ text "Priority"
                , input [ id "priority", type_ "text" ] []
                ]
            , div []
                [ text "Status"
                , input [ id "status", type_ "text" ] []
                ]
            , div []
                [ button [ type_ "submit" ]
                    [ text "Add to list" ]
                ]
            ]
        ]

main : Html msg
main =
    view initialModel