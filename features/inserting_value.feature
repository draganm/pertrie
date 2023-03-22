Feature: inserting value


    Scenario: inserting one value in an empty database
        Given an empty database
        When I insert a value into the root trie
        Then the size of the root trie should be 1
