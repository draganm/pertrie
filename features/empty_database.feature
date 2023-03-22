Feature: empty database

    Scenario: zero size of the empty database
        Given an empty database
        When I get the size of the root
        Then the size should be 0
