# Strawpoll Plugin

## Details

* This plugin will allow you to perform one of the following actions
  * Get the status of a poll by it's ID
  * Create a poll
    * To create a poll there needs to be a title and at least two options
  * Options
    * Array of options that are used for the poll
  * Multi allows for someone to vote more than once
    * True or False
  * Dup is how to determine if the user has voted or not
    * Normal: IP based
    * Permissive: cookie based
    * Disabled: No check
  * Captcha determines if a captcha is placed on the poll
    * True or false

```
CMD: /strawpoll
Help: /strawpoll {id | title [options] (multi) (dup) (captcha)}
```

### Get Poll Example

```
/strawpoll 15832795
---------------
Title: This is a test poll.
URL: https://www.strawpoll.me/15832795
Option: Vote Count
Option #1: 0
Option #2: 0
```

### Make Poll Example

```
/strawpoll "This is my title" "something, something2, something3" "false" "normal" "true"
---------------
Title: This is my title
URL: https://www.strawpoll.me/15833955
```