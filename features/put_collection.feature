Feature: Put Collection

  Scenario: PUT /collections
    Given I have these collections:
            """
            [
                {
                    "id": "00112233-4455-6677-8899-aabbccddeeff",
                    "e_tag": "45678",
                    "name": "Coronavirus key indicators",
                    "publish_date": "2020-05-10T14:58:29.317Z"
                }
            ]
            """
    When I set the "If-Match" header to "45678"
    And I PUT "/collections/00112233-4455-6677-8899-aabbccddeeff"
            """
            {
                "name": "Coronavirus key indicators",
                "publish_date": "2020-05-05T14:58:29.317Z"
            }
            """
    Then the HTTP status code should be "200"

  Scenario: PUT /collections with an out of date ETag
    Given I have these collections:
            """
            [
                {
                    "id": "00112233-4455-6677-8899-aabbccddeeff",
                    "e_tag": "45678",
                    "name": "Coronavirus key indicators",
                    "publish_date": "2020-05-10T14:58:29.317Z"
                }
            ]
            """
    When I set the "If-Match" header to "1111"
    And I PUT "/collections/00112233-4455-6677-8899-aabbccddeeff"
            """
            {
                "name": "Coronavirus key indicators",
                "publish_date": "2020-05-05T14:58:29.317Z"
            }
            """
    Then the HTTP status code should be "409"
    And I should receive the following JSON response:
        """
        {
            "errors":[ {"message":  "out of date collection resource"}]
        }
        """

  Scenario: PUT /collections with non existent collection
    Given I have these collections:
            """
            []
            """
    When I set the "If-Match" header to "45678"
    And I PUT "/collections/00112233-4455-6677-8899-aabbccddeeff"
            """
            {
                "name": "Coronavirus key indicators",
                "publish_date": "2020-05-05T14:58:29.317Z"
            }
            """
    Then the HTTP status code should be "404"
    And I should receive the following JSON response:
        """
        {
            "errors":[ {"message":  "collection not found"}]
        }
        """