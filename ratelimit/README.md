# Rate Limit Plugins

This is a background plugin and does not have a command, or populate the help menu.

## Rate Limit

* The default time between commands a user can send is 2.5 seconds
    * If a user breaks this rule the timer gets reset on their account
* The default broken rules limit is 10
    * If a user breaks the rule limit their time between commands grows by a second
    * Every time they break the rule an extra second gets added to their wait time.