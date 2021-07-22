Feature: Update collection
  Scenario: Update the publish date of a collection
    Given I have these collections:
        """
            [
                {
                    "id": "00112233-4455-6677-8899-aabbccddeeff",
                    "name": "LMSV1",
                    "publish_date": "2020-05-10T14:58:29.317Z"
                },
                {
                    "id": "10112233-4455-6677-8899-aabbccddeeff",
                    "name": "LMSV2",
                    "publish_date": "2020-05-05T14:58:29.317Z"
                },
                {
                    "id": "20112233-4455-6677-8899-aabbccddeeff",
                    "name": "LMSV3",
                    "publish_date": "2020-05-08T14:58:29.317Z"
                }
            ]
        """
    When I PUT "/collections/00112233-4455-6677-8899-aabbccddeeff"
    """
        {
          "id": "00112233-4455-6677-8899-aabbccddeeff",
          "name": "LMSV1",
          "publish_date": "2020-06-10T14:58:29.317Z"
        }
    """
    Then the HTTP status code should be "200"
