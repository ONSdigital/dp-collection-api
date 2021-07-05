Feature: Post Collection
    Scenario: POST /collections
        Given there are no collections
        When I POST "/collections"
            """
            {
                "name": "Coronavirus key indicators",
                "publish_date": "2020-05-05T14:58:29.317Z"
            }
            """
        Then the HTTP status code should be "201"